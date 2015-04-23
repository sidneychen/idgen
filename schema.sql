CREATE TABLE `id_gen` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `service` varchar(255) NOT NULL DEFAULT '',
  `position` int(10) unsigned NOT NULL,
  `update_time` int(10) unsigned NOT NULL,
  `status` tinyint(3) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_service` (`service`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;