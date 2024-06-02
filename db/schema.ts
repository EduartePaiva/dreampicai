import {
    serial,
    text,
    timestamp,
    pgTable,
    uuid,
    pgSchema,
    boolean,
    integer,
} from "drizzle-orm/pg-core";

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

export const image = pgTable("images", {
    id: uuid("id").notNull().defaultRandom().primaryKey(),
    userId: uuid("user_id")
        .references(() => users.id)
        .notNull(),
    imageLocation: text("image_location"),
    status: integer("status").notNull().default(1),
    prompt: text("prompt").notNull(),
    deleted: boolean("deleted").notNull().default(false),
    createdAt: timestamp("created_at").defaultNow().notNull(),
    deletedAt: timestamp("deleted_at"),
});

type user = typeof user.$inferSelect;
type image = typeof image.$inferSelect;
