CREATE TABLE `comments` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `post_id` int(10) unsigned NOT NULL,
    `name` varchar(255) NOT NULL,
    `email` varchar(255) NOT NULL,
    `body` text NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `posts` (
     `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
     `user_id` int(10) unsigned NOT NULL,
     `title` varchar(255) NOT NULL,
     `body` text NOT NULL,
     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

