import { ZResponse } from '@go-kickstart/zod'
import { match } from 'ts-pattern'

export const getSecurityMetadata = ({
	security = true,
	securityType = 'bearer',
}: {
	security?: boolean
	securityType?: 'bearer' | 'service'
} = {}) => {
	const openApiSecurity = match(securityType)
		.with('bearer', () => [
			{
				bearerAuth: [],
			},
		])
		.with('service', () => [
			{
				'x-service-token': [],
			},
		])
		.exhaustive()

	return {
		...(security && { openApiSecurity }),
	}
}

export const failResponses = {
	401: ZResponse,
	403: ZResponse,
	404: ZResponse,
	500: ZResponse,
} as const
