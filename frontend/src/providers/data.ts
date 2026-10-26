import type { DataProvider, HttpError } from "@refinedev/core"
import { createSimpleRestDataProvider } from "@refinedev/rest/simple-rest"
import { call, ApiError } from "@/api/client"
import type {
  CatalogModel,
  CatalogModelCreateInput,
  CatalogModelList,
  CatalogModelUpdateInput,
  User,
  UserCreateInput,
  UserList,
  UserUpdateInput,
  VirtualKey,
  VirtualKeyCreateInput,
  VirtualKeyList,
  VirtualKeyUpdateInput,
} from "@/api/generated/contract"
import {
  createModel,
  createUser,
  createVirtualKey,
  deleteModel,
  deleteUser,
  deleteVirtualKey,
  getModel,
  getUser,
  getVirtualKey,
  updateModel,
  updateUser,
  updateVirtualKey,
} from "@/api/generated/operations"
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
  getList: async (params) => {
    if (params.resource === "users") {
      try {
        const currentPage = params.pagination?.currentPage ?? 1
        const pageSize = params.pagination?.pageSize ?? 10
        const start = (currentPage - 1) * pageSize
        const end = start + pageSize
        const list = await call<UserList>({
          method: "GET",
          path: `/api/users?_start=${start}&_end=${end}`,
        })
        return { data: list.items, total: list.total }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "models") {
      try {
        const currentPage = params.pagination?.currentPage ?? 1
        const pageSize = params.pagination?.pageSize ?? 10
        const start = (currentPage - 1) * pageSize
        const end = start + pageSize
        const list = await call<CatalogModelList>({
          method: "GET",
          path: `/api/models?_start=${start}&_end=${end}`,
        })
        return { data: list.items, total: list.total }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "virtual-keys") {
      try {
        const currentPage = params.pagination?.currentPage ?? 1
        const pageSize = params.pagination?.pageSize ?? 10
        const start = (currentPage - 1) * pageSize
        const end = start + pageSize
        const list = await call<VirtualKeyList>({
          method: "GET",
          path: `/api/virtual-keys?_start=${start}&_end=${end}`,
        })
        return { data: list.items, total: list.total }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    return restProvider.getList!(params)
  },
  getOne: async (params) => {
    if (params.resource === "profile") {
      try {
        return { data: await me() }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "users") {
      try {
        return { data: await getUser({ id: Number(params.id) }) }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "models") {
      try {
        return { data: await getModel({ id: Number(params.id) }) }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "virtual-keys") {
      try {
        return { data: await getVirtualKey({ id: Number(params.id) }) }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    return restProvider.getOne!(params)
  },
  create: async (params) => {
    if (params.resource === "users") {
      try {
        const input = params.variables as UserCreateInput
        const data = await createUser(input)
        return { data }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "models") {
      try {
        const input = params.variables as CatalogModelCreateInput
        const data = await createModel({
          name: input.name,
          displayName: input.displayName,
          type: input.type,
          disabled: Boolean(input.disabled),
        })
        return { data }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "virtual-keys") {
      try {
        const input = params.variables as VirtualKeyCreateInput
        const data = await createVirtualKey({
          name: input.name,
          maxBudget: input.maxBudget ?? "",
          expiresIn: input.expiresIn ?? "",
          budgetReset: input.budgetReset ?? "",
        })
        return { data }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    return restProvider.create!(params)
  },
  update: async (params) => {
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
    if (params.resource === "users") {
      try {
        const input = params.variables as Omit<UserUpdateInput, "id">
        const data = await updateUser({
          id: Number(params.id),
          name: input.name,
          email: input.email,
          password: input.password ?? "",
          role: input.role,
          disabled: Boolean(input.disabled),
        })
        return { data }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "models") {
      try {
        const input = params.variables as Omit<CatalogModelUpdateInput, "id">
        const data = await updateModel({
          id: Number(params.id),
          name: input.name,
          displayName: input.displayName,
          type: input.type,
          disabled: Boolean(input.disabled),
        })
        return { data }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "virtual-keys") {
      try {
        const input = params.variables as Omit<VirtualKeyUpdateInput, "id">
        const data = await updateVirtualKey({
          id: Number(params.id),
          name: input.name,
          disabled: Boolean(input.disabled),
          maxBudget: input.maxBudget ?? "",
          expiresIn: input.expiresIn ?? "",
          budgetReset: input.budgetReset ?? "",
        })
        return { data }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    return restProvider.update!(params)
  },
  deleteOne: async (params) => {
    if (params.resource === "users") {
      try {
        await deleteUser({ id: Number(params.id) })
        return { data: { id: Number(params.id) } as User }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "models") {
      try {
        await deleteModel({ id: Number(params.id) })
        return { data: { id: Number(params.id) } as CatalogModel }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    if (params.resource === "virtual-keys") {
      try {
        await deleteVirtualKey({ id: Number(params.id) })
        return { data: { id: Number(params.id) } as VirtualKey }
      } catch (error) {
        return Promise.reject(toHttpError(error))
      }
    }
    return restProvider.deleteOne!(params)
  },
} as DataProvider

export { kyInstance }
