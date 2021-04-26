import { useSelector } from 'react-redux'
import { useEffect, useState } from 'react'

import Footer from '../components/Footer'
import Header from '../components/Header'
import { Redirect, useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next';
import GroupService from '../services/group';

const Group = () => {
  const { shareHash } = useParams()
  const { t } = useTranslation(['translation', 'profile']);
  // eslint-disable-next-line no-unused-vars
  const { token, role } = useSelector(
    (state) => state.auth
  )
  const [redirectStarted, setRedirectStarted] = useState(false)

  useEffect(()=> {
    GroupService.acceptGroupShare(shareHash)
    .then(
      (response) => {
        setTimeout(() => {
            setRedirectStarted(true)
        }, 2000);
      },
      (error) => {
        console.log(error.message)
      }
    )
  }, [shareHash])

  if (!token) {
    return <Redirect to='/login' />
  }

  if (redirectStarted) {
    return <Redirect to='/profile' />
  }
  
  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center pb-4'>
          <div className='pt-4 text-xl font-bold text-purple-900'>
            {t('profile:groupShareAccept', 'Accept Group Share')}
          </div>
        </div>
        <div className='p-4 bg-purple-300 border border-purple-300 rounded w-1-2'>
            <p>{t('profile:groupShareStarted', 'Accepting group share...')}</p>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default Group