DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM pg_type WHERE typname = 'sensor_type') THEN
        ALTER TYPE sensor_type ADD VALUE 'CAR_DIRECTION';
        ALTER TYPE sensor_type ADD VALUE 'FRONT_ULTRASONIC';
    END IF;
END $$;