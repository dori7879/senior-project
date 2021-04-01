import { useTranslation } from 'react-i18next';
import { useDispatch } from 'react-redux'
import { useForm } from "react-hook-form";
import { useState } from "react";
import { register } from '../actions/auth'



const RegistrationForm = () => {
  const { t } = useTranslation(['translation', 'registration']);
  const { register, handleSubmit } = useForm();
  const dispatch = useDispatch()
  const [successful, setSuccesfull] = useState(false)
  const onSubmit = data => {
    console.log(data);
    setSuccesfull(true)
  }

  return (
    <form className='mt-8'onSubmit={handleSubmit(onSubmit)} name='registration'>
      {successful ? (
        <div className='flex items-center justify-center'>
        <h2 className='mt-3 mb-10 text-2xl font-extrabold leading-7 text-center text-purple-900 uppercase'>
          {t('registration:registered', 'Registered')}
        </h2>
      </div>
      ) : (
          <div>
          <div>
              <h2 className='mt-6 text-2xl font-extrabold leading-7 text-center text-purple-900 uppercase'>
                {t('registration:caption', 'Register')}
              </h2>
          </div>
          <input type='hidden' name='remember' value='true' />
          <div className='rounded shadow-sm'>
            <div>
              <label className='text-sm text-gray-500'>{t('registration:firstname', 'First Name')}</label>
              <input
                name='FirstName'
                type='firstname'
                required
                ref={register({ required: true, maxLength: 30 })}
                className='relative block w-full px-3 py-2 text-gray-800 placeholder-purple-400 border border-purple-200 rounded appearance-none focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5'
              />
            </div>
            <div>
              <label className='text-sm text-gray-500'>{t('registration:lastname', 'Last Name')}</label>
              <input
                name='LastName'
                type='lastname'
                required
                ref={register({ required: true, maxLength: 30 })}
                className='relative block w-full px-3 py-2 text-gray-800 placeholder-purple-400 border border-purple-200 rounded appearance-none focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5'
              />
            </div>
            <div className='mt-2 mb-1'>
              <label className='mr-3 text-sm text-gray-500'>{t('registration:role', 'Role')}</label>
              <select
                name='IsTeacher'
                ref={register}
                className='text-xs bg-purple-100 border border-purple-300'
              >
                <option value={true}>{t('registration:teacher', 'Teacher')}</option>
                <option value={false}>{t('registration:student', 'Student')}</option>
              </select>
            </div>
            <div>
              <label className='text-sm text-gray-500'>{t('registration:email', 'Email address')}</label>
              <input
                name='Email'
                type='email'
                required
                ref={register({ required: true, maxLength: 30 })}
                className='relative block w-full px-3 py-2 text-gray-800 placeholder-purple-400 border border-purple-200 rounded appearance-none focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5'
              />
            </div>
            <div>
              <label className='text-sm text-gray-500'>{t('registration:password', 'Password')}</label>
              <input
                name='Password'
                type='password'
                required
                ref={register({ required: true, maxLength: 30 })}
                className='relative block w-full px-3 py-2 text-gray-800 placeholder-purple-400 border border-purple-200 rounded appearance-none focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5'
              />
            </div>
          </div>
          <div className='pb-4 mt-4'>
            <button
              type='submit'
              className='relative flex justify-center w-full px-2 py-2 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
            >
              {t('registration:signup', 'Sign Up')}
            </button>
          </div>
        </div>
       ) }

    </form>
  )
}

export default RegistrationForm
