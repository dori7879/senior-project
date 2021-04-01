import { useDispatch, useSelector } from 'react-redux'

import { Redirect } from 'react-router-dom'
import { login } from '../actions/auth'
import { useForm } from "react-hook-form";
import { useTranslation } from 'react-i18next';

const required = (value) => {
    if (!value) {
      return (
        <div
          className='text-center text-red-800 bg-red-200 border border-rounded'
          role='alert'
        >
          This field is required!
        </div>
      )
    }
  }
const LoginForm = () => {
  const { t } = useTranslation(['translation', 'login']);
  const isLoggedIn = useSelector((state) => state.auth.isLoggedIn)
  const { register, handleSubmit, errors } = useForm();
  const onSubmit = data =>  console.log(data);
 // const dispatch = useDispatch()


  if (isLoggedIn) return <Redirect to='/profile' />

  return (
    <form name='loginForm' onSubmit={handleSubmit(onSubmit)} className='mt-8' ref={register}>
      
      <div className='rounded shadow-sm'>
        <div>
          <input
            name='Email'
            required
            type='email'
            className='relative block w-full px-3 py-2 text-gray-800 placeholder-purple-400 border border-purple-200 rounded appearance-none focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5'
            placeholder={t('login:email', 'Email address')}
            ref={register({ required: true, maxLength: 30 })}
          />
            {errors.name && errors.name.type === "required" && (
             <span role="alert">This is required</span>
            )}
            {errors.name && errors.name.type === "maxLength" && (
                <span role="alert">Max length exceeded</span>
            )}
        </div>
        <div className='mt-4'>
          <input
            name='Password'
            type='password'
            required
            className='relative block w-full px-3 py-2 text-gray-800 placeholder-purple-400 border border-purple-200 rounded appearance-none focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5'
            placeholder={t('login:password', 'Password')}
            ref={register({ required: true, maxLength: 30 })}
          />
        </div>
      </div>

      <div className='flex flex-row justify-between mt-2'>
        <div className='flex items-center'>
          <input name='remember-me' type='checkbox' />
          <span className='text-xs text-purple-900'>{t('login:remember', 'Remember me')}</span>
        </div>
        <div className='flex items-center'>
          <a
            href='/'
            className='text-xs font-medium text-purple-900 transition duration-150 ease-in-out hover:text-purple-700 focus:outline-none focus:underline'
          >
            {t('login:forgotpass', 'Forgot your password?')}
          </a>
        </div>
      </div>
      <div className='pb-3 mt-3'>
        <button
          type='submit'
          className='relative flex justify-center w-full px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
        >
          {t('login:signin', 'Sign in')}
        </button>
      </div>
      
    </form>
  )
}

export default LoginForm
