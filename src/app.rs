use clap::Parser;
use colored::*;
use dirs;
use rusqlite::{params, Connection};
use std::fmt::Display;
use std::fs::OpenOptions;
use std::panic;
use std::path::Path;

#[derive(Parser)]
#[command(name = "iktodo", version = "0.0.1", author = "alctny")]
#[command(about = "A simple todo list CLI app")]
enum Command {
    /// 新增代办事项
    Add {
        /// 开始时间
        #[arg(short, long)]
        begin: Option<String>,
        /// 结束时间
        #[arg(short, long)]
        end: Option<String>,
        /// 任务名称
        task: Vec<String>,
    },
    /// 查看所有代办事项
    List,
    /// 修改代办事项
    Modify {
        /// 任务 ID
        id: i32,
        /// 新的任务名称
        name: String,
    },
    /// 完成代办事项
    Done {
        /// 任务ID
        id: Vec<i32>,
    },
    /// 清除所有已完成事项、过期事项
    Clean,
    /// 初始化数据库
    Init,
}

#[derive(Parser)]
struct ListFlag {
    /// 通过标签分组展示
    #[arg(short, long)]
    bytag: Option<Option<i32>>,
}

pub struct Task {
    id: i32,
    name: String,
    status: i32, // 任务状态：0-未完成，1-已完成
    create_at: String,
    begin_at: Option<String>,
    finish_at: Option<String>,
    // _repeat: Option<u8>, // 重复周期
}

impl Task {
    fn get_begin_at(&self) -> String {
        match self.begin_at {
            Some(ref begin_at) => begin_at.to_string(),
            None => self.create_at.to_string(),
        }
    }

    fn get_end_at(&self) -> String {
        match self.finish_at {
            Some(ref end_at) => end_at.to_string(),
            None => "".to_string(),
        }
    }
}

impl Display for Task {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let msg = format!("{:02} {}", self.id, self.name);
        match self.status {
            0 => write!(f, "{}", msg.green()),
            1 => write!(f, "\x1B[9m{}\x1B[0m", msg.yellow()),
            2 => write!(f, "{}", msg.red()),
            _ => write!(f, "{}", msg),
        }
    }
}

pub fn exec() {
    let app = Command::parse();
    match app {
        Command::Add { begin, end, task } => add_task(&Task {
            id: 0,
            name: task.join(" "),
            status: 3,
            create_at: "".to_string(),
            begin_at: begin,
            finish_at: end,
        }),
        Command::List => list_task(),
        Command::Modify { id, name } => modify_task(id, &name),
        Command::Done { id } => done_task(id),
        Command::Clean => clean_task(),
        Command::Init => init_db(),
    }
}

fn init_db() {
    let homedir = match dirs::home_dir() {
        Some(path) => path,
        None => Path::new(".").to_path_buf(),
    };
    let dbpath = homedir.join(".config/iktodo.db");
    if !dbpath.exists() {
        OpenOptions::new()
            .write(true)
            .create(true)
            .open(&dbpath)
            .expect("create db file error");
    }

    let conn = match Connection::open(&dbpath) {
        Ok(conn) => conn,
        Err(e) => {
            println!("Error: {}", e);
            return;
        }
    };

    let status = conn.execute(
        "CREATE TABLE IF NOT EXISTS task (
            id 		    INTEGER     PRIMARY KEY,
            name 	    TEXT,
            status	    INTEGER     NOT NULL DEFAULT 0,
            create_at   DATE        NOT NULL DEFAULT CURRENT_DATE,
            begin_at    DATE,
            finish_at   DATE
        );
        ",
        [],
    );

    match status {
        Ok(_) => println!("init db success"),
        Err(e) => println!("Error to init db:\n {}", e),
    }
}

// 连接数据库
fn connect_sqlite() -> Connection {
    let dbpath = match dirs::home_dir() {
        Some(path) => path.join(".config/iktodo.db"),
        None => panic!("can't find home dir"),
    };
    if !Path::new(&dbpath).exists() {
        panic!("db file not exists, you can use `iktodo init` to create it.")
    }
    Connection::open(&dbpath).expect("connection db failed")
}

// 新增代办事项
fn add_task(t: &Task) {
    let conn = connect_sqlite();
    let sql = match (&t.begin_at, &t.finish_at) {
        (None, None) => format!("INSERT INTO task (name) VALUES ('{}');", t.name),
        (Some(b), None) => format!(
            "INSERT INTO task (name, begin_at) VALUES ('{}', '{}');",
            t.name, b
        ),
        (None, Some(f)) => format!(
            "INSERT INTO task (name, finish_at) VALUES ('{}', '{}');",
            t.name, f
        ),
        (Some(b), Some(f)) => format!(
            "INSERT INTO task (name, begin_at, finish_at) VALUES ('{}', '{}', '{}');",
            t.name, b, f
        ),
    };
    let status = conn.execute(&sql, []);
    match status {
        Ok(_) => list_task(),
        Err(e) => panic!("Error in add task: {}", e),
    }
}

// 查看所有代办事项
fn list_task() {
    let sql = "SELECT id, name, status, create_at, begin_at, finish_at FROM task ORDER BY status ASC, begin_at ASC, finish_at ASC, create_at ASC;";
    let conn = connect_sqlite();
    let mut stmt = conn.prepare(sql).ok().unwrap();
    let task_list = stmt.query_map(params![], |row| {
        Ok(Task {
            id: row.get(0)?,
            name: row.get(1)?,
            status: row.get(2)?,
            create_at: row.get(3)?,
            begin_at: row.get(4)?,
            finish_at: row.get(5)?,
        })
    });
    let task_list = match task_list {
        Ok(list) => list,
        Err(e) => panic!("Error in query list task: {}", e),
    };

    for task_result in task_list {
        let task = task_result.unwrap();
        println!("{}", task);
    }
}

// 将代办事项标记为已完成
fn done_task(id: Vec<i32>) {
    let conn = connect_sqlite();
    let ids = id
        .iter()
        .map(|x| x.to_string())
        .collect::<Vec<String>>()
        .join(", ");
        
    let stat = conn.prepare(&format!("UPDATE task SET status = 1 WHERE id IN ({})", ids));
    let mut stat = match stat {
        Ok(stat) => stat,
        Err(e) => panic!("Error in done task: {}", e),
    };
    match stat.execute([]) {
        Ok(_) => list_task(),
        Err(e) => println!("Error in mark task: {}", e),
    }
}

// 修改代办事项
fn modify_task(id: i32, name: &str) {
    let conn = connect_sqlite();
    let status = conn.execute(
        "UPDATE task SET name = ?1 WHERE id = ?2",
        &[name, &id.to_string()],
    );
    match status {
        Ok(_) => list_task(),
        Err(e) => panic!("Error in modify task: {}", e),
    }
}

// 清理已完成代办事项
fn clean_task() {
    let conn = connect_sqlite();
    let status = conn.execute("DELETE FROM task WHERE status = 1", []);
    match status {
        Ok(_) => list_task(),
        Err(e) => panic!("Error in clean task: {}", e),
    }
}
