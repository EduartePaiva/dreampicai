import { serial, text, timestamp, pgTable, uuid, pgSchema } from "drizzle-orm/pg-core";

const authSchema = pgSchema("auth");
const users = authSchema.table("users", {
    id: uuid("id").primaryKey(),
});

export const user = pgTable("accounts", {
    id: serial("id").primaryKey(),
    userId: uuid("user_id").references(() => users.id),
    userName: text("user_name").notNull(),
    createdAt: timestamp("created_at").defaultNow().notNull(),
});

type user = typeof user.$inferSelect;
