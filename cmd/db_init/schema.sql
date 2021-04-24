/* Users Table */
CREATE TABLE IF NOT EXISTS users
(
 "id"         bigserial NOT NULL,
 email        text NOT NULL,
 name         text NOT NULL,
 display_name text NOT NULL,
 auth_key     text NOT NULL,
 pwd_hash     text NULL,
 pwd_salt     text NULL,
 created_at   date NOT NULL,
 updated_at   date,
 CONSTRAINT PK_users PRIMARY KEY ( "id" )
);

/* Projects Table */
CREATE TABLE IF NOT EXISTS projects
(
 "id"           bigserial NOT NULL,
 created_by   int NOT NULL,
 display_name text NOT NULL,
 description  text NULL,
 created_at   date NOT NULL,
 updated_at   date,
 CONSTRAINT PK_projects PRIMARY KEY ( "id" ),
 CONSTRAINT FK_34 FOREIGN KEY ( created_by ) REFERENCES users ( "id" )
);

CREATE INDEX IF NOT EXISTS fkIdx_35 ON projects
(
 created_by
);

/* Devices Table */
CREATE TABLE IF NOT EXISTS devices
(
 "id"         bigserial NOT NULL,
 project_id   int NOT NULL,
 display_name text NOT NULL,
 auth_key     text NOT NULL UNIQUE,
 description  text NULL,
 created_at   date NOT NULL,
 updated_at   date,
 CONSTRAINT PK_devices PRIMARY KEY ( "id" ),
 CONSTRAINT FK_37 FOREIGN KEY ( project_id ) REFERENCES projects ( "id" )
);

CREATE INDEX IF NOT EXISTS fkIdx_38 ON devices
(
 project_id
);

/* Endpoints Table */
CREATE TABLE IF NOT EXISTS endpoints
(
 "id"           int NOT NULL,
 device_id    int NOT NULL,
 display_name text NOT NULL,
 description  text NULL,
 "pattern"      text NOT NULL,
 created_at   date NOT NULL,
 updated_at   date,
 CONSTRAINT PK_endpoints PRIMARY KEY ( "id" ),
 CONSTRAINT FK_48 FOREIGN KEY ( device_id ) REFERENCES devices ( "id" )
);

CREATE INDEX IF NOT EXISTS fkIdx_49 ON endpoints
(
 device_id
);

/* External Keys Table */
CREATE TABLE IF NOT EXISTS external_keys
(
 "id"         int NOT NULL,
 auth_key   text NOT NULL,
 created_at date NOT NULL,
 updated_at date NOT NULL,
 created_by int NOT NULL,
 CONSTRAINT PK_external_keys PRIMARY KEY ( "id" ),
 CONSTRAINT FK_61 FOREIGN KEY ( created_by ) REFERENCES users ( "id" )
);

CREATE INDEX IF NOT EXISTS fkIdx_62 ON external_keys
(
 created_by
);
