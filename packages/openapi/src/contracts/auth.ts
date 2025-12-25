import {
	ZAuthGoogleLoginDTO,
	ZAuthLoginDTO,
	ZAuthRegisterDTO,
	ZAuthResult,
	ZAuthVerifyEmailDTO,
	ZAuthVerifyEmailResponse,
} from '@go-kickstart/zod'
import { initContract } from '@ts-rest/core'
import { failResponses } from '../utils.js'

const c = initContract()

export const authContract = c.router({
	register: {
		summary: 'Register',
		description: 'Register a new user',
		path: '/api/v1/auth/register',
		method: 'POST',
		body: ZAuthRegisterDTO,
		responses: {
			201: ZAuthResult,
			...failResponses,
		},
	},
	login: {
		summary: 'Login',
		description: 'Login with email/username and password',
		path: '/api/v1/auth/login',
		method: 'POST',
		body: ZAuthLoginDTO,
		responses: {
			200: ZAuthResult,
			...failResponses,
		},
	},
	googleLogin: {
		summary: 'Google login',
		description: 'Login or register using Google OAuth',
		path: '/api/v1/auth/google',
		method: 'POST',
		body: ZAuthGoogleLoginDTO,
		responses: {
			200: ZAuthResult,
			...failResponses,
		},
	},
	verifyEmail: {
		summary: 'Verify email',
		description: 'Verify user email using a verification code',
		path: '/api/v1/auth/verify-email',
		method: 'POST',
		body: ZAuthVerifyEmailDTO,
		responses: {
			200: ZAuthVerifyEmailResponse,
			...failResponses,
		},
	},
})
