import { z } from 'zod'

export const ZResponse = z.object({
	status: z.number().int(),
	message: z.string(),
	success: z.boolean(),
})

export function ZResponseWithData<T>(schema: z.ZodSchema<T>) {
	return z.object({ data: schema }).extend(ZResponse.shape)
}

export function ZPaginatedResponse<T>(schema: z.ZodSchema<T>) {
	return z
		.object({
			total: z.number(),
			page: z.number(),
			limit: z.number(),
			totalPages: z.number(),
			data: z.array(schema),
		})
		.extend(ZResponse.shape)
}

export const ZModel = z.object({
	id: z.string().uuid(),
	createdAt: z.string().datetime(),
	updatedAt: z.string().datetime(),
	deletedAt: z.string().datetime().optional(),
})

export const ZGetManyQuery = z.object({
	limit: z.coerce.number().int().nonnegative().optional(),
	offset: z.coerce.number().int().nonnegative().optional(),
	preloads: z.string().optional(),
	orderBy: z.string().optional(),
	orderDirection: z.enum(['asc', 'desc']).optional(),
})

export const ZPreloadsQuery = z.object({
	preloads: z.string().optional(),
})
