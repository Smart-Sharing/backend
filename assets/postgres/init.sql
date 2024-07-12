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
  password varchar (128) NOT NUll,

	CHECK (job_position IN ('worker', 'admin')),
  CHECK (LENGTH(password) >= 8),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sessions(
  id SERIAL,
  state integer DEFAULT 0,
  machine_id varchar(16) NOT NULL,
  worker_id integer NOT NULL,
  datetime_start TIMESTAMP DEFAULT CURRENT_DATE,
  datetime_finish TIMESTAMP DEFAULT CURRENT_DATE,

  CHECK (state IN (0, 1, 2)),

  PRIMARY KEY (id),

  FOREIGN KEY (machine_id) REFERENCES machines (id) ON DELETE CASCADE,
  FOREIGN KEY (worker_id) REFERENCES users (id) ON DELETE CASCADE
);

-- CREATE TEST DATA into database
INSERT INTO users(name, phone_number, job_position, password) 
  VALUES ('USER1', '89099769897', 'worker', '12345678'), 
         ('SUPER-USER', '89090001122', 'admin', 'some-password123');

INSERT INTO machines(id) VALUES ('1FGH345'), ('1ASD987');

INSERT INTO sessions(machine_id, worker_id) VALUES ('1FGH345', 2), ('1ASD987', 1);


