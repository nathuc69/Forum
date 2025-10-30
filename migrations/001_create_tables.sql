PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password TEXT,
	email TEXT,
	token TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS topics (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	user_id INTEGER NOT NULL,
	UNIQUE(title, content, user_id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	topic_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	UNIQUE(content, topic_id, user_id),
	FOREIGN KEY (topic_id) REFERENCES topics(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS reactions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	value INTEGER NOT NULL CHECK(value IN (-1, 1)), /* -1 = dislike, 1 = like */
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	target_type TEXT NOT NULL CHECK(target_type IN ('topics','posts')), /*liste des "cibles": objets sur lesquels mettre des réactions*/
	target_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	UNIQUE(target_type, target_id, user_id), /* un unique vote par utilisateur par cible */
	FOREIGN KEY (user_id) REFERENCES users(id)
);

  CREATE TABLE IF NOT EXISTS topic_categories (
      topic_id INT REFERENCES topics(id) ON DELETE CASCADE,
      category_id INT REFERENCES categories(id) ON DELETE CASCADE,
      PRIMARY KEY (topic_id, category_id)
  );