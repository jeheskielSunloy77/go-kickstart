import { initContract } from '@ts-rest/core'

import {
	ZAuthGoogleLoginRequest,
	ZAuthLoginRequest,
	ZAuthRegisterRequest,
	ZAuthResult,
	ZAuthVerifyEmailRequest,
	ZAuthVerifyEmailResponse,
} from '@go-kickstart/zod'

const c = initContract()

export const authContract = c.router({
	register: {
		summary: 'Register',
		path: '/api/v1/auth/register',
		method: 'POST',
		body: ZAuthRegisterRequest,
		responses: {
			201: ZAuthResult,
		},
	},
	login: {
		summary: 'Login',
		path: '/api/v1/auth/login',
		method: 'POST',
		body: ZAuthLoginRequest,
		responses: {
			200: ZAuthResult,
		},
	},
	googleLogin: {
		summary: 'Google login',
		path: '/api/v1/auth/google',
		method: 'POST',
		body: ZAuthGoogleLoginRequest,
		responses: {
			200: ZAuthResult,
		},
	},
	verifyEmail: {
		summary: 'Verify email',
		path: '/api/v1/auth/verify-email',
		method: 'POST',
		body: ZAuthVerifyEmailRequest,
		responses: {
			200: ZAuthVerifyEmailResponse,
		},
	},
})
