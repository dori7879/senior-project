import { useDispatch, useSelector } from 'react-redux'
import { useRef, useState } from 'react'

import Footer from '../components/Footer'
import Form from 'react-validation/build/form'
import Header from '../components/Header'
import Input from 'react-validation/build/input'
import { Redirect } from 'react-router-dom'
import authHeader from '../services/auth-header'
import axios from 'axios'
import { useTranslation } from 'react-i18next';


const ChangePassword = () => {
    const { t } = useTranslation(['translation', 'profile']);
    const [oldPassword, setOldPassword] = useState('')
    const [newPassword, setNewPassword] = useState('')
    const { access_token, role } = useSelector(
      (state) => state.auth
    )
    const [success, setSuccess] = useState(false)
  
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
  
    const handleSubmit = (e) => {
      e.preventDefault()
  
      axios
        .post(
          'https://radiant-inlet-12251.herokuapp.com/api/v1/reset-password',
          {
            old_password: oldPassword,
            new_password: newPassword,
          },
          {
            headers: authHeader(),
          }
        )
        .then(() => {
          setSuccess(true)
        })
        .catch((er) => {
          //history.push('/signin')
          console.log(er.message)
        })
    }

    return(
        <div>
          <p className='mt-3 text-purple-900'>
            <strong>Change Password </strong>
          </p>
          <Form className='mt-8' onSubmit={handleSubmit}>
            <div className='flex flex-row items-center pb-2 items'>
              <div className='flex flex-col'>
                <label className='block w-full px-4 pt-1 mb-2 ml-5 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                  {t('profile:old', 'Old Password')}
                </label>
                <Input
                  className='px-2 py-1 mr-6 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                  aria-label='Password'
                  value={oldPassword}
                  onChange={(e) => setOldPassword(e.target.value)}
                  validations={[required]}
                  name='oldPassword'
                  type='password'
                  required
                  placeholder={t('profile:password', 'Password')}
                />
              </div>
              <div className='flex flex-col'>
                <label className='block w-full px-4 pt-1 mb-2 ml-3 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                  {t('profile:new', 'New Password')}
                </label>
                <Input
                  className='px-2 py-1 mr-6 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                  aria-label='Password'
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  validations={[required]}
                  name='newPassword'
                  type='password'
                  required
                  placeholder={t('profile:password', 'Password')}
                />
              </div>
            </div>
            <button
              type='submit'
              className='relative flex justify-center w-full px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
            >
              {t('profile:save', 'Save')}
            </button>
          </Form>
          {success ? (
            <div className='p-4 bg-purple-300 border border-purple-300 rounded w-1-2'>
              <p className='text-purple-900'>
                <strong>{t('profile:change', 'Password has changed')} </strong>
              </p>
            </div>
          ) : null}  
        </div>
        
    )

}

export default ChangePassword;