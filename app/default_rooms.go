package app

type DefaultRoomType struct {
	id          string
	name        string
	cover       string
	slug        string
	description string
}

func GetDefaultRooms() []DefaultRoomType {
	defaultRoomList := []DefaultRoomType{
		{"1", "General Support", "general.jpg",
			"general-support",
			`This channel provides a safe and inclusive space for
		programmers to seek and offer general support. Whether you're feeling overwhelmed, need someone
		to talk to, or want to share your experiences, this channel is here for you. Together, we can
		navigate the challenges and support one another in our journeys.`},
		{"2", "Self-Care and Well-being", "happy.jpg", "self-care-and-well-being", `Taking care of ourselves is crucial, especially
		when dealing with depression and pressure. In this channel, we focus on discussing and sharing
		self-care tips, strategies for managing stress, mindfulness techniques, and ways to promote
		overall well-being. Join us to explore various approaches and learn how to prioritize your
		mental health.`},
		{"3", "Career and Burnout", "depressed.jpg",
			"career-and-burnout",
			`Balancing a demanding career in programming can
		sometimes lead to burnout. This channel is dedicated to discussing career-related challenges,
		work-life balance, imposter syndrome, and strategies to prevent burnout. Share your experiences,
		seek advice, and discover ways to build a sustainable and fulfilling career in programming.`},
		{"4", "Positive Vibes", "playful.jpg",
			"positive-vibes",
			`Sometimes, we all need a little positivity and inspiration.
		In this channel, let's focus on uplifting and motivating each other. Share success stories,
		celebrate achievements, and spread positive vibes. Together, we can create a supportive
		environment that reminds us of our strengths and encourages us to keep moving forward.`},
	}
	return defaultRoomList
}
