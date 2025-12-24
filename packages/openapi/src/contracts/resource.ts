import { initContract } from '@ts-rest/core'
import { z } from 'zod'

import { schemaWithPagination } from '@go-kickstart/zod'
import { getSecurityMetadata } from '../utils.js'

type ResourceContractOptions = {
	path: string
	resource: string
	resourcePlural: string
	schemas: {
		entity: z.ZodTypeAny
		create: z.ZodTypeAny
		update: z.ZodTypeAny
	}
	security?: boolean
	securityType?: 'bearer' | 'service'
}

const c = initContract()

const userIdParams = z.object({
	id: z.string().uuid(),
})

const preloadsQuery = z.object({
	preloads: z.string().optional(),
})

const listQuery = z.object({
	limit: z.coerce.number().int().nonnegative().optional(),
	offset: z.coerce.number().int().nonnegative().optional(),
	preloads: z.string().optional(),
	order_by: z.string().optional(),
	order_direction: z.enum(['asc', 'desc']).optional(),
})

export const createResourceContract = ({
	path,
	resource,
	resourcePlural,
	schemas,
	security = true,
	securityType = 'bearer',
}: ResourceContractOptions) => {
	const metadata = getSecurityMetadata({ security, securityType })

	return c.router({
		getMany: {
			summary: `Get Many ${resourcePlural}`,
			description: `Retrieve a paginated list of ${resourcePlural} that can be filtered, sorted, and preloaded.`,
			path,
			method: 'GET',
			query: listQuery,
			responses: {
				200: schemaWithPagination(schemas.entity),
			},
			metadata,
		},
		getById: {
			summary: `Get ${resource} by ID`,
			description: `Retrieve a single ${resource} by its unique identifier (ID), with optional preloaded relationships.`,
			path: `${path}/:id`,
			method: 'GET',
			pathParams: userIdParams,
			query: preloadsQuery,
			responses: {
				200: schemas.entity,
			},
			metadata,
		},
		store: {
			summary: `Store ${resource}`,
			description: `Create a new ${resource} with the provided data, with validation and will return the created entity.`,
			path,
			method: 'POST',
			body: schemas.create,
			responses: {
				201: schemas.entity,
			},
			metadata,
		},
		update: {
			summary: `Update ${resource}`,
			description: `Update an existing ${resource} by its ID with the provided data, and return the updated entity.`,
			path: `${path}/:id`,
			method: 'PATCH',
			pathParams: userIdParams,
			body: schemas.update,
			responses: {
				200: schemas.entity,
			},
			metadata,
		},
		destroy: {
			summary: `Destroy ${resource}`,
			description: `Soft delete the specified ${resource} by its ID. This action is reversible.`,
			path: `${path}/:id`,
			method: 'DELETE',
			pathParams: userIdParams,
			responses: {
				204: c.noBody(),
			},
			metadata,
		},
		kill: {
			summary: `Kill ${resource}`,
			description: `Permanently delete the specified ${resource} by its ID. This action is irreversible.`,
			path: `${path}/:id/kill`,
			method: 'DELETE',
			pathParams: userIdParams,
			responses: {
				204: c.noBody(),
			},
			metadata,
		},
		restore: {
			summary: `Restore ${resource}`,
			description: `Restore a previously soft-deleted ${resource} by its ID and then return the restored entity.`,
			path: `${path}/:id/restore`,
			method: 'PATCH',
			pathParams: userIdParams,
			query: preloadsQuery,
			body: c.noBody(),
			responses: {
				200: schemas.entity,
			},
			metadata,
		},
	})
}
