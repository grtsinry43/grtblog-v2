import type { PostDetail } from "$lib/features/post/types";
import { createModelDataContext } from "svatoms";

export const postDetailCtx = createModelDataContext<PostDetail | null>({
    name: "postDetailCtx",
    initial: null,
});