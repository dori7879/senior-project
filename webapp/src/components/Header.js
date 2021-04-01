import { useDispatch, useSelector } from 'react-redux'

import { Link } from 'react-router-dom'
import { logout } from '../actions/auth'
import { useTranslation } from 'react-i18next';

  
const Header = () => {
const isLoggedIn = useSelector((state) => state.auth.isLoggedIn)

  const dispatch = useDispatch()
  const { t, i18n } = useTranslation(['translation', 'header']);

  const changeLanguage = code => {
    i18n.changeLanguage(code);
  };
  return (
    <div className='flex justify-center w-full px-4 py-2 bg-purple-900'>
      <div className='w-full max-w-3xl '>
        <div className='flex items-center justify-between text-purple-200'>
          <Link to='/'>
            <div className='font-mono text-2xl text-purple-100 cursor-pointer select-none hover:text-purple-300'>
              EasySubmit
            </div>
          </Link>
          
          <div className='flex'>
          <button className='text-center w-8 h-8 mr-2 text-s font-thin border-2 border-purple-100 rounded-lg focus:outline-none active:border-purple-900 hover:text-purple-300 hover:border-purple-300 ' type="button" onClick={() => changeLanguage('ru')}>ru</button>
          <button className='text-center w-8 h-8  mr-2 text-s font-thin border-2 border-purple-100 rounded-lg focus:outline-none active:border-purple-900 hover:text-purple-300 hover:border-purple-300 ' type="button" onClick={() => changeLanguage('en')}>en</button>
            {isLoggedIn ? (
              <div>
                <Link to='/profile'>
                  <button
                    type='button'
                    className='text-center w-auto h-8 p-1 mr-2 text-xs font-thin border-2 border-purple-100 rounded-lg focus:outline-none active:border-purple-900 hover:text-purple-300 hover:border-purple-300 '
                  >
                    {t('header:profile', 'My profile')}
                  </button>
                </Link>
                <Link to='/'>
                  <button
                    type='button'
                    onClick={() => dispatch(logout)}
                    className='text-center w-auto h-8 p-1 text-xs font-thin border-2 border-purple-100 rounded-lg focus:outline-none hover:text-purple-300 hover:border-purple-300'
                  >
                    {t('header:signout', 'Sign out')}
                  </button>
                </Link>
              </div>
            ) : (
              <div>
                <Link to='/signin'>
                  <button
                    type='button'
                    className='text-center w-auto h-8 p-1 mr-2 text-xs font-thin border-2 border-purple-100 rounded-lg focus:outline-none active:border-purple-900 hover:text-purple-300 hover:border-purple-300 '
                  >
                    {t('header:signin', 'Sign in')}
                  </button>
                </Link>
                <Link to='/signup'>
                  <button
                    type='button'
                    className=' h-8 p-1 text-center w-auto text-xs font-thin border-2 border-purple-100 rounded-lg focus:outline-none hover:text-purple-300 hover:border-purple-300'
                  >
                    {t('header:signup', 'Sign up')}
                  </button>
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default Header
