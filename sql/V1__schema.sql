CREATE TABLE `pack_sizes` (
    `size` INT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `pack_sizes` (`size`) VALUES (250), (500), (1000), (2000), (5000);
