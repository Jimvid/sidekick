import { Bounce, toast } from 'react-toastify'
import type { ToastOptions } from 'react-toastify'

const baseOptions: ToastOptions = {
  position: 'top-center',
  autoClose: 3000,
  hideProgressBar: false,
  closeOnClick: false,
  pauseOnHover: true,
  draggable: true,
  progress: undefined,
  theme: 'dark',
  transition: Bounce,
}
export const notifySuccess = (message: string, options?: ToastOptions) =>
  toast.success(message, {
    ...baseOptions,
    ...options,
  })

export const notifyWarning = (message: string, options?: ToastOptions) =>
  toast.warning(message, {
    ...baseOptions,
    ...options,
  })

export const notifyError = (message: string, options?: ToastOptions) =>
  toast.error(message, {
    ...baseOptions,
    ...options,
  })

export const notifyInfo = (message: string, options?: ToastOptions) =>
  toast.info(message, {
    ...baseOptions,
    ...options,
  })
