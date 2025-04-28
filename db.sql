CREATE DATABASE atm_demo;

USE atm_demo;

CREATE TABLE accounts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    pin INT NOT NULL,
    balance DOUBLE DEFAULT 0
);
