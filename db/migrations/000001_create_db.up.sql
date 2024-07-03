CREATE TABLE IF NOT EXISTS skill (
    skill_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS skill_group (
    skill_group_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS skill_group_skill (
    skill_group_id INTEGER NOT NULL,
    skill_id INTEGER NOT NULL,
    PRIMARY KEY (skill_group_id, skill_id),
    FOREIGN KEY (skill_group_id) REFERENCES skill_group (skill_group_id),
    FOREIGN KEY (skill_id) REFERENCES skill (skill_id)
);