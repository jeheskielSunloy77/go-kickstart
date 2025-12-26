import { ZResponse } from '@go-kickstart/zod'

export const getSecurityMetadata = ({
	security = true,
	securityType = 'bearer',
}: {
	security?: boolean
	securityType?: 'bearer' | 'service'
} = {}) => {
	const openApiSecurity = (() => {
		switch (securityType) {
			case 'bearer':
				return [{ bearerAuth: [] }]
			case 'service':
				return [{ 'x-service-token': [] }]
			default:
				const _exhaustive: never = securityType
				throw new Error(`Unhandled securityType: ${_exhaustive}`)
		}
	})()

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
