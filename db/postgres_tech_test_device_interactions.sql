CREATE DATABASE IF NOT EXISTS tech_test;
\c tech_test;

-- Table structure for table device_interactions
DROP TABLE IF EXISTS device_interactions;
CREATE TABLE device_interactions (
  interaction_id SERIAL PRIMARY KEY,
  timestamp TIMESTAMPTZ NOT NULL,
  latitude DECIMAL(10,8) NOT NULL,
  longitude DECIMAL(11,8) NOT NULL,
  device_id INT NOT NULL,
  device_name VARCHAR(255) NOT NULL
);

-- Dumping data for table device_interactions
-- Your data insertion statements go here
