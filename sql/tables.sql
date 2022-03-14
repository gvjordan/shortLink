START TRANSACTION;

CREATE TABLE IF NOT EXISTS `links` (
  `ID` int(10) NOT NULL AUTO_INCREMENT,
  `Name` varchar(10) NOT NULL,
  `URL` text NOT NULL,
  `CreatedAt` int(11) UNSIGNED NOT NULL,
  `CreatedBy` text NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `tokens` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `token` varchar(25) NOT NULL,
  `description` varchar(255) NOT NULL,
  `ip` varchar(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

COMMIT;
