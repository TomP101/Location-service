DROP TABLE IF EXISTS user;
CREATE TABLE user (
  id         INT AUTO_INCREMENT NOT NULL,
  name      VARCHAR(128) NOT NULL,
  latitude     DECIMAL(11,8) NOT NULL,
  longitude      DECIMAL(11,8) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO user
  (name, latitude, longitude)
VALUES
  ('Blue Train', 13.123, 25.123),
  ('Giant Steps', 14.123, 24.123),
  ('Jeru', 15.123, 23.123),
  ('Sarah Vaughan', 16.123, 22.123);