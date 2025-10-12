import { handleApiError as he, type ApiError, formatDate } from '@yuelioi/utils'

export function handleApiError(err: unknown, action: string) {
  he(err, {
    action: action,
    onError: (err: ApiError) => toast.error(err.message),
  })
}

export { formatDate }
