import { defineConfig } from "drizzle-kit";
import "dotenv/config";

export default defineConfig({
    schemaFilter: ["public"],
    schema: "./db/schema.ts",
    dialect: "postgresql",
    dbCredentials: {
        url: process.env.SUPABASE_DRIZZLE_URL!,
    },
    verbose: true,
    strict: true,
    out: "./cmd/migrate",
});
