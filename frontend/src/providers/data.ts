import type { DataProvider, HttpError } from "@refinedev/core"
import { createSimpleRestDataProvider } from "@refinedev/rest/simple-rest"
import { ApiError } from "@/api/client"
import { me, updateProfile, type UpdateProfileInput } from "@/api/session"
import { API_URL } from "./constants"

const { dataProvider: restProvider, kyInstance } =
  createSimpleRestDataProvider({
    apiURL: API_URL,
    kyOptions: {
      credentials: "include",
    },
  })

function toHttpError(error: unknown): HttpError {
  if (error instanceof ApiError) {
    return {
      message: error.message,
      statusCode: error.status,
      errors: error.fields,
    }
  }
  if (error instanceof Error) {
    return { message: error.message, statusCode: 500 }
  }
  return { message: "request failed", statusCode: 500 }
}

export const dataProvider = {
  ...restProvider,
  getOne: async (params: Parameters<NonNullable<DataProvider["getOne"]>>[0]) => {
    if (params.resource === "profile") {
      try {
        return { data: await me() }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    return restProvider.getOne!(params)
  },
  update: async (
    params: Parameters<NonNullable<DataProvider["update"]>>[0]
  ) => {
    if (params.resource === "profile") {
      try {
        const input = params.variables as UpdateProfileInput
        return {
          data: await updateProfile({
            name: input.name,
            email: input.email,
            currentPassword: input.currentPassword,
            newPassword: input.newPassword,
          }),
        }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    return restProvider.update!(params)
  },
} as DataProvider

export { kyInstance }
