import logo from '../assets/fav.ico'
import { useTranslation } from 'react-i18next';

const Footer = () => {
    const { t } = useTranslation(['translation', 'footer']);
    return (
      <div className='bg-purple-200'>
        <footer className='text-gray-800 body-font'>
          <div className='container flex flex-col flex-wrap px-24 py-12 mx-auto md:items-center lg:items-start md:flex-row md:flex-no-wrap'>
            <div className='flex-shrink-0 w-32 mx-auto mt-10 text-center md:mx-0 md:text-left md:mt-0 '>
              <a href='/' className='flex items-center justify-center font-medium text-gray-900 title-font md:justify-start'>
                <img className='w-8 h-8' src={logo} alt='app logo' />
                <span className='ml-1 font-bold text-purple-900 text-md'>
                  EasySubmit
                </span>
              </a>
              <p className='mt-1 ml-10 text-xs text-justify text-gray-600'>
               {t('footer:description', 'A simple and robust system that makes the interactions as intuitive as possible.')}
              </p>
            </div>
            <div className='flex flex-wrap flex-grow order-first -mb-10 text-center md:pr-20 md:text-left'>
              <div className='w-full px-4 lg:w-1/3 md:w-1/2'>
                <h2 className='mb-3 text-sm font-medium tracking-widest text-purple-900 title-font'>
                  {t('footer:about' , 'ABOUT US')}
                </h2>
                <nav className='mb-10 list-none'>
                  <li>
                    <a href='/' className='text-sm text-gray-600 hover:text-purple-700 hover:underline'>
                      {t('footer:company', 'Company')}
                    </a>
                  </li>
                  <li>
                    <a href='/' className='text-sm text-gray-600 hover:text-purple-700 hover:underline'>
                      {t('footer:news', 'News')}
                    </a>
                  </li>
                  <li>
                    <a href='/' className='text-sm text-gray-600 hover:text-purple-700 hover:underline'>
                      {t('footer:policies', 'Policies')}
                    </a>
                  </li>
                </nav>
              </div>
              <div className='w-full px-4 lg:w-1/3 md:w-1/2'>
                <h2 className='mb-3 text-sm font-medium tracking-widest text-purple-900 title-font'>
                  {t('footer:followus', 'FOLLOW US')}
                </h2>
                <nav className='mb-10 list-none'>
                  <li>
                    <a href='/' className='text-sm text-gray-600 hover:text-purple-700 hover:underline'>
                      Facebook
                    </a>
                  </li>
                  <li>
                    <a href='/' className='text-sm text-gray-600 hover:text-purple-700 hover:underline'>
                      Instagram
                    </a>
                  </li>
                  <li>
                    <a href='/' className='text-sm text-gray-600 hover:text-purple-700 hover:underline'>
                      Twitter
                    </a>
                  </li>
                </nav>
              </div>
              <div className='w-full pl-4 lg:w-1/3 md:w-1/2'>
                <h2 className='mb-3 text-sm font-medium tracking-widest text-purple-900 title-font'>
                  {t('footer:suppport', 'SUPPORT')}
                </h2>
                <nav className='mb-10 list-none'>
                  <li>
                    <a  href='/' className='text-sm text-gray-600 hover:text-purple-700 hover:underline'>
                      {t('footer:faq', 'FAQ')}
                    </a>
                  </li>
                  <li>
                    <a href='/' className='text-sm text-gray-600 hover:text-purple-700 hover:underline'>
                      {t('footer:contacts', 'Contacts')}
                    </a>
                  </li>
                </nav>
              </div>
            </div>
          </div>
          <div className='container flex flex-wrap justify-center w-screen max-w-full px-5 py-2 mx-auto bg-purple-100 '>
            <p className='text-xs text-gray-600 '>
             {t('footer:rights', ' © 2020 All Rights Reserved. EasySubmit')}
            </p>
          </div>
        </footer>
      </div>
    )
  }
  
  export default Footer