ALTER TABLE
  ONLY public.car_sessions
ADD
  CONSTRAINT fk_car_session_car FOREIGN KEY (car_id) REFERENCES public.cars(id);

ALTER TABLE
  ONLY public.car_controllers
ADD
  CONSTRAINT fk_car_controller_car FOREIGN KEY (car_id) REFERENCES public.cars(id);

ALTER TABLE
  ONLY public.car_controllers
ADD
  CONSTRAINT fk_car_controller_controller_instance FOREIGN KEY (controller_instance_id) REFERENCES public.controller_instances(id);

ALTER TABLE
  ONLY public.car_session_controllers
ADD
  CONSTRAINT fk_car_session_controller_car_session FOREIGN KEY (car_session_id) REFERENCES public.car_sessions(id);

ALTER TABLE
  ONLY public.car_session_controllers
ADD
  CONSTRAINT fk_car_session_controller_controller_instance FOREIGN KEY (controller_instance_id) REFERENCES public.controller_instances(id);