CREATE TABLE `climate` (
    `id` int NOT NULL AUTO_INCREMENT,
    `datetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `temperature1` double NOT NULL,
    `humidity1` double NOT NULL,
    PRIMARY KEY (`id`)
)