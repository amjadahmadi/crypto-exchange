db.createUser(
    {
        user: "mongoadmin",
        pwd: "bdung",
        roles: [
            {
                role: "readWrite",
                db: "trade"
            }
        ]
    }
);