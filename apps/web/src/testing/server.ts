import { setupServer } from "msw/node";
import { handlers } from "@/testing/handlers";

export const server = setupServer(...handlers);
