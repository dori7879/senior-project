import { useDispatch, useSelector } from 'react-redux'
import { useRef, useState } from 'react'

import Footer from '../components/Footer'
import Header from '../components/Header'
import { Redirect } from 'react-router-dom'

import { useTranslation } from 'react-i18next';
import ChangePassword from '../components/ChangePassword'
import PersonalData from '../components/PersonalData'

const Profile = () => {
  const { t } = useTranslation(['translation', 'profile']);
  const { access_token, role } = useSelector(
    (state) => state.auth
  )

  if (!access_token) {
    //return <Redirect to='/signin' />
  }
  
  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center pb-4'>
          <div className='pt-4 text-xl font-bold text-purple-900'>
            {t('profile:caption', 'My Profile')}
          </div>
        </div>
        <div className='p-4 bg-purple-300 border border-purple-300 rounded w-1-2'>
          <PersonalData />
          <ChangePassword />
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default Profile