import { useEffect } from "react"
import authHeader from '../services/auth-header'
import axios from 'axios'
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux'

const PersonalData = () => {
      const { t } = useTranslation(['translation', 'profile']); 
      const FirstName='fdsfd'
      const LastName = 'fsdfgsdsd'
      const Email = 'sdfdsgfdgfd'
      const { AccessToken, IsTeacher} = useSelector(
        (state) => state.auth
      )
    useEffect(()=> {
     
    })
    return(
        <div>
            <p className='text-purple-900'>
            <strong>{t('profile:firstname', 'First Name:')} </strong> {FirstName}
          </p>
          <p className='text-purple-900'>
            <strong>{t('profile:lastname', 'Last Name:')}  </strong> {LastName}
          </p>
          <p className='text-purple-900'>
            <strong>{t('profile:email', 'Email:')} </strong> {Email}
          </p>
          <p className='text-purple-900'>
            <strong>{t('profile:role', 'Role:')} </strong> {IsTeacher ? t('profile:teacher', 'Teacher') : t('profile:student', 'Student')}
          </p>
          
        </div>
    )
}

export default PersonalData;