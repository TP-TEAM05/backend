CREATE TYPE sensor_type AS ENUM (
  'FRONT_LIDAR',
  'FRONT_ULTRASONIC',
  'REAR_ULTRASONIC',
  'WHEEL_SPEED',
  'GPS_LOCATION',
  'GPS_SPEED',
  'GPS_DIRECTION',
  'MAGNETOMETER_DIRECTION'
);

-- CREATE TABLES
CREATE TABLE public.car_controllers (
  id bigserial NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  car_id bigint NOT NULL,
  controller_instance_id bigint NOT NULL
);

CREATE TABLE public.car_session_controllers (
  id bigserial NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  car_session_id bigint,
  controller_instance_id bigint
);

CREATE TABLE public.car_sessions (
  id bigserial NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  car_id bigint,
  session_id bigint,
  is_controlled_by_user boolean
);

CREATE TABLE public.cars (
  id bigserial NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  vin text,
  name text,
  color text
);

CREATE TABLE public.controller_instances (
  id bigserial NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  firmware_id bigint,
  controller_id bigint
);

CREATE TABLE public.controllers (
  id bigserial NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  name character varying(255),
  type character varying(255),
  description text
);

CREATE TABLE public.firmwares (
  id bigserial NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  version character varying(255),
  description text
);

CREATE TABLE public.sensors (
  id bigserial NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  controller_instance_id bigint,
  name character varying(255),
  sensor_type sensor_type
);

CREATE TABLE public.sessions (
  id bigserial NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  name character varying(255),
  started_at timestamp with time zone,
  ended_at timestamp with time zone
);

CREATE TABLE public.measurements (
  car_session_id bigint,
  sensor_id bigint,
  created_at timestamp with time zone,
  data1 double precision,
  data2 double precision
);

-- CREATE HYPER_TABLES
SELECT
  create_hypertable('measurements', 'created_at');

-- ADD CONSTRAINTS
ALTER TABLE
  ONLY public.car_controllers
ADD
  CONSTRAINT car_controllers_pkey PRIMARY KEY (id);

ALTER TABLE
  ONLY public.car_session_controllers
ADD
  CONSTRAINT car_session_controllers_pkey PRIMARY KEY (id);

ALTER TABLE
  ONLY public.car_sessions
ADD
  CONSTRAINT car_sessions_pkey PRIMARY KEY (id);

ALTER TABLE
  ONLY public.cars
ADD
  CONSTRAINT cars_pkey PRIMARY KEY (id);

ALTER TABLE
  ONLY public.controller_instances
ADD
  CONSTRAINT controller_instances_pkey PRIMARY KEY (id);

ALTER TABLE
  ONLY public.controllers
ADD
  CONSTRAINT controllers_pkey PRIMARY KEY (id);

ALTER TABLE
  ONLY public.firmwares
ADD
  CONSTRAINT firmwares_pkey PRIMARY KEY (id);

ALTER TABLE
  ONLY public.sensors
ADD
  CONSTRAINT sensors_pkey PRIMARY KEY (id);

ALTER TABLE
  ONLY public.sessions
ADD
  CONSTRAINT sessions_pkey PRIMARY KEY (id);

ALTER TABLE
  public.measurements
ADD
  CONSTRAINT fk_car_sessions_measurements FOREIGN KEY (car_session_id) REFERENCES public.car_sessions(id);

ALTER TABLE
  ONLY public.car_sessions
ADD
  CONSTRAINT fk_car_sessions_session FOREIGN KEY (session_id) REFERENCES public.sessions(id);

ALTER TABLE
  ONLY public.controller_instances
ADD
  CONSTRAINT fk_controllers_controller_instances FOREIGN KEY (controller_id) REFERENCES public.controllers(id);

ALTER TABLE
  ONLY public.controller_instances
ADD
  CONSTRAINT fk_firmwares_controller_instances FOREIGN KEY (firmware_id) REFERENCES public.firmwares(id);

ALTER TABLE
  public.measurements
ADD
  CONSTRAINT fk_measurements_sensor FOREIGN KEY (sensor_id) REFERENCES public.sensors(id);

ALTER TABLE
  ONLY public.sensors
ADD
  CONSTRAINT fk_sensors_controller_instance FOREIGN KEY (controller_instance_id) REFERENCES public.controller_instances(id);

-- CREATE INDEXES
CREATE INDEX idx_car_controllers_deleted_at ON public.car_controllers USING btree (deleted_at);

CREATE INDEX idx_car_session_controllers_deleted_at ON public.car_session_controllers USING btree (deleted_at);

CREATE INDEX idx_car_sessions_deleted_at ON public.car_sessions USING btree (deleted_at);

CREATE INDEX idx_cars_deleted_at ON public.cars USING btree (deleted_at);

CREATE INDEX idx_controller_instances_deleted_at ON public.controller_instances USING btree (deleted_at);

CREATE INDEX idx_controllers_deleted_at ON public.controllers USING btree (deleted_at);

CREATE INDEX idx_firmwares_deleted_at ON public.firmwares USING btree (deleted_at);

CREATE INDEX idx_sensors_deleted_at ON public.sensors USING btree (deleted_at);

CREATE INDEX idx_sessions_deleted_at ON public.sessions USING btree (deleted_at);