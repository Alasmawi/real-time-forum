package database

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"reboot01.com/js/realtime-forum/internal/security"
)

func (db DB) SeedData() error {

	// Hash the hashed_password `123` using bcrypt
	hashedPassword, err := security.Hash(`123`)
	if err != nil {
		return err
	}

	// insertUsers := []string{
	// 	`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_hashed_password, session_id)
	// 	VALUES ('John', 'Doe', 25, 1, 'johndoe', 'john.doe@example.com', '` + hashedPassword + `', NULL);`,
	// }

	insertCategories := []string{
		`INSERT INTO category (name) VALUES ('Technology')`,
		`INSERT INTO category (name) VALUES ('Sports')`,
		`INSERT INTO category (name) VALUES ('Music')`,
		`INSERT INTO category (name) VALUES ('Movies')`,
		`INSERT INTO category (name) VALUES ('Gaming')`,
		`INSERT INTO category (name) VALUES ('Food')`,
		`INSERT INTO category (name) VALUES ('Travel')`,
		`INSERT INTO category (name) VALUES ('Science')`,
		`INSERT INTO category (name) VALUES ('Art')`,
		`INSERT INTO category (name) VALUES ('Books')`,
	}

	insertUsers := []string{
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password, session_id) VALUES ('John', 'Doe', 25, 1, 'johndoe', 'john.doe@example.com', '` + hashedPassword + `', NULL);`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Jane', 'Smith', 28, 0, 'janesmith', 'jane.smith@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Mike', 'Johnson', 32, 1, 'mikejohnson', 'mike.johnson@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Sarah', 'Williams', 24, 0, 'sarahwilliams', 'sarah.williams@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('David', 'Brown', 30, 1, 'davidbrown', 'david.brown@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Emily', 'Davis', 26, 0, 'emilydavis', 'emily.davis@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Chris', 'Wilson', 29, 1, 'chriswilson', 'chris.wilson@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Lisa', 'Moore', 27, 0, 'lisamoore', 'lisa.moore@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Tom', 'Taylor', 31, 1, 'tomtaylor', 'tom.taylor@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Anna', 'Anderson', 23, 0, 'annaanderson', 'anna.anderson@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Mark', 'Thomas', 35, 1, 'markthomas', 'mark.thomas@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Jessica', 'Jackson', 22, 0, 'jessicajackson', 'jessica.jackson@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Ryan', 'White', 33, 1, 'ryanwhite', 'ryan.white@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Amanda', 'Harris', 28, 0, 'amandaharris', 'amanda.harris@email.com', '` + hashedPassword + `', NULL)`,
		`INSERT INTO user (f_name, l_name, age, sex, username, email, hashed_password) VALUES ('Kevin', 'Martin', 26, 1, 'kevinmartin', 'kevin.martin@email.com', '` + hashedPassword + `', NULL)`,
	}

	insertPosts := []string{
		`INSERT INTO post (content, created_at, user_id) VALUES ('Just discovered this amazing new JavaScript framework! It makes building UIs so much easier. Has anyone else tried it?', '2024-01-15 09:30:00', 1)`,
		`INSERT INTO post (content, created_at, user_id) VALUES ('Watching the game tonight! My team is finally in the playoffs after 5 years. So excited!', '2024-01-15 18:45:00', 2)`,
		`INSERT INTO post (content, created_at, user_id) VALUES ('Been listening to this new album on repeat. The guitar work is absolutely incredible. Highly recommend!', '2024-01-16 14:20:00', 3)`,
		`INSERT INTO post (content, created_at, user_id) VALUES ('Just finished watching the latest superhero movie. The special effects were mind-blowing! What did everyone think?', '2024-01-16 21:15:00', 4)`,
		`INSERT INTO post (content, created_at, user_id) VALUES ('Finally beat that boss level that had me stuck for weeks! The satisfaction is real. Gaming win!', '2024-01-17 11:00:00', 5)`,
		`INSERT INTO post (content, created_at, user_id) VALUES ('Tried making homemade pasta for the first time. It was surprisingly therapeutic and delicious! Recipe in comments.', '2024-01-17 19:30:00', 6)`,
		`INSERT INTO post (content, created_at, user_id) VALUES ('Planning my summer vacation to Japan. Any recommendations for must-visit places in Tokyo?', '2024-01-18 08:15:00', 7)`,
		`INSERT INTO post (content, created_at, user_id) VALUES ('Read an fascinating article about quantum computing today. The future of technology is amazing!', '2024-01-18 16:45:00', 8)`,
		`INSERT INTO post (content, created_at, user_id) VALUES ('Finished my first oil painting! It took months but Im so proud of the result. Art is truly rewarding.', '2024-01-19 13:00:00', 9)`,
		`INSERT INTO post (content, created_at, user_id) VALUES ('Just finished reading this incredible fantasy novel. The world-building was phenomenal! Anyone else read it?', '2024-01-19 20:30:00', 10)`,
	}

	insertPostCategories := []string{
		`INSERT INTO post_has_category (post_id, category_id) VALUES (1, 1)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (1, 8)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (1, 10)`,

		`INSERT INTO post_has_category (post_id, category_id) VALUES (2, 2)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (2, 7)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (2, 6)`,

		`INSERT INTO post_has_category (post_id, category_id) VALUES (3, 3)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (3, 9)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (3, 4)`,

		`INSERT INTO post_has_category (post_id, category_id) VALUES (4, 4)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (4, 1)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (4, 5)`,

		`INSERT INTO post_has_category (post_id, category_id) VALUES (5, 5)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (5, 1)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (5, 2)`,

		`INSERT INTO post_has_category (post_id, category_id) VALUES (6, 6)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (6, 7)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (6, 9)`,

		`INSERT INTO post_has_category (post_id, category_id) VALUES (7, 7)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (7, 6)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (7, 10)`,

		`INSERT INTO post_has_category (post_id, category_id) VALUES (8, 8)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (8, 1)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (8, 10)`,

		`INSERT INTO post_has_category (post_id, category_id) VALUES (9, 9)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (9, 3)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (9, 7)`,

		`INSERT INTO post_has_category (post_id, category_id) VALUES (10, 10)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (10, 4)`,
		`INSERT INTO post_has_category (post_id, category_id) VALUES (10, 8)`,
	}

	insertComments := []string{

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('I tried it last week and its fantastic! The learning curve is pretty gentle too.', '2024-01-15 10:15:00', 1, 11)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Thanks for sharing! Ive been looking for something like this for my project.', '2024-01-15 11:30:00', 1, 12)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('The documentation is really well written, which makes all the difference.', '2024-01-15 12:45:00', 1, 13)`,

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Go team! Ive been following them since I was a kid. This is their year!', '2024-01-15 19:00:00', 2, 14)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('I have tickets to the game! The atmosphere is going to be electric.', '2024-01-15 19:30:00', 2, 15)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Their defense has been solid all season. Should be a great match!', '2024-01-15 20:00:00', 2, 1)`,

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Yes! The third track is my absolute favorite. Such intricate melodies.', '2024-01-16 15:00:00', 3, 2)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Been a fan of this artist for years. They never disappoint!', '2024-01-16 15:45:00', 3, 3)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Added to my playlist immediately. Thanks for the recommendation!', '2024-01-16 16:30:00', 3, 4)`,

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('The action sequences were incredible! Best movie of the year so far.', '2024-01-16 22:00:00', 4, 5)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('I loved the character development. Much better than the previous one.', '2024-01-16 22:30:00', 4, 6)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('The plot twist caught me completely off guard. Brilliant writing!', '2024-01-16 23:00:00', 4, 7)`,

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Congrats! That boss took me forever too. What strategy did you use?', '2024-01-17 11:30:00', 5, 8)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('The feeling of finally beating a tough boss is unmatched in gaming!', '2024-01-17 12:00:00', 5, 9)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Now you can move on to the even harder levels! Good luck!', '2024-01-17 12:30:00', 5, 10)`,

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Homemade pasta is the best! Would love to see that recipe.', '2024-01-17 20:00:00', 6, 11)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('I tried making pasta once and it was a disaster. Any tips for beginners?', '2024-01-17 20:30:00', 6, 12)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('The key is getting the dough consistency just right. Practice makes perfect!', '2024-01-17 21:00:00', 6, 13)`,

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Tokyo is amazing! Make sure to visit the fish market early in the morning.', '2024-01-18 09:00:00', 7, 14)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Shibuya crossing is a must-see, and the food scene is incredible everywhere.', '2024-01-18 09:30:00', 7, 15)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('I recommend getting a JR Pass for easy travel around the city and country.', '2024-01-18 10:00:00', 7, 1)`,

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Quantum computing is going to revolutionize so many fields! Exciting times.', '2024-01-18 17:15:00', 8, 2)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('The implications for cryptography alone are mind-boggling.', '2024-01-18 17:45:00', 8, 3)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('I read that article too! The potential applications seem almost limitless.', '2024-01-18 18:15:00', 8, 4)`,

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Oil painting is so rewarding but challenging. Your patience paid off!', '2024-01-19 13:30:00', 9, 5)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Id love to see a photo of your painting! Art takes such dedication.', '2024-01-19 14:00:00', 9, 6)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('Thinking about trying painting myself. Any advice for a complete beginner?', '2024-01-19 14:30:00', 9, 7)`,

		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('That book is on my reading list! Your review just moved it to the top.', '2024-01-19 21:00:00', 10, 8)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('The authors world-building skills are legendary. Everything feels so real.', '2024-01-19 21:30:00', 10, 9)`,
		`INSERT INTO comment (content, created_at, post_id, user_id) VALUES ('I finished the whole series in a week! Couldnt put them down.', '2024-01-19 22:00:00', 10, 10)`,
	}

	allInserts := [][]string{
		insertCategories,
		insertUsers,
		insertPosts,
		insertPostCategories,
		insertComments,
	}

	for _, group := range allInserts {
		for _, stmt := range group {
			_, err := db.Exec(stmt)
			if err != nil {
				return fmt.Errorf(`failed to execute insert statement: %w`, err)
			}
		}
	}

	fmt.Println(`Database seeded successfully`)
	return nil
}
