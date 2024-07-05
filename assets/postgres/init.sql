CREATE TABLE IF NOT EXISTS machines(
	id varchar(16) NOT NULL,
	state integer DEFAULT 0,
  voltage integer DEFAULT 0,
		
	CHECK (state IN (0, 1)),
  CHECK (voltage >= 0),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users(
  id SERIAL,
	name text NOT NULL,
	phone_number varchar(11) NOT NULL UNIQUE,
	job_position varchar(8) NOT NULL,

	CHECK (job_position IN ('worker', 'operator', 'admin')),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sessions(
  id SERIAL,
  state integer DEFAULT 0,
  machine_id varchar(16) NOT NULL,
  worker_id integer NOT NULL,
  datetime_start TIMESTAMP,
  datetime_finish TIMESTAMP,

  CHECK (state IN (0, 1, 2)),

  PRIMARY KEY (id),

  FOREIGN KEY (machine_id) REFERENCES machines (id) ON DELETE CASCADE,
  FOREIGN KEY (worker_id) REFERENCES users (id) ON DELETE CASCADE
);

