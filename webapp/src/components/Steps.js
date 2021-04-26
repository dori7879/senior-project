import { useTranslation } from 'react-i18next';

const Steps = () => {
    const { t } = useTranslation(['translation', 'steps']);
    return (
      <div>
        <div className='flex justify-center pt-4 text-xl font-bold text-purple-900 border-t-2 border-purple-200'>
          {t('steps:caption', 'How to use EasySubmit?')}
        </div>
        <div className='flex justify-center'>
          <section className='text-gray-600 body-font'>
            <div className='container flex flex-wrap px-5 py-4 mx-auto'>
              <div className='relative flex pt-10 pb-20 mx-auto sm:items-center md:w-2/3'>
                <div className='absolute inset-0 flex items-center justify-center w-6 h-full'>
                  <div className='w-1 h-full bg-purple-200 pointer-events-none'></div>
                </div>
                <div className='relative z-10 inline-flex items-center justify-center flex-shrink-0 w-6 h-6 mt-8 text-sm font-medium text-purple-100 bg-purple-800 rounded-full sm:mt-0 title-font'>
                  1
                </div>
                <div className='flex flex-col items-start flex-grow pl-6 md:pl-8 sm:items-center sm:flex-row'>
                  <div className='flex-grow mt-6 sm:pl-6 sm:mt-0'>
                    <h2 className='mb-1 text-xl font-medium text-purple-700 title-font'>
                      {t('steps:step1', 'Choose type of the assignment')}
                    </h2>
                    <p className='leading-relaxed'>
                      {t('steps:desc1', 'Create a homework (quiz or attendance) page by clicking on the corresponding button on the Home page')}
                    </p>
                  </div>
                </div>
              </div>
              <div className='relative flex pb-20 mx-auto sm:items-center md:w-2/3'>
                <div className='absolute inset-0 flex items-center justify-center w-6 h-full'>
                  <div className='w-1 h-full bg-purple-200 pointer-events-none'></div>
                </div>
                <div className='relative z-10 inline-flex items-center justify-center flex-shrink-0 w-6 h-6 mt-10 text-sm font-medium text-purple-100 bg-purple-800 rounded-full sm:mt-0 title-font'>
                  2
                </div>
                <div className='flex flex-col items-start flex-grow pl-6 md:pl-8 sm:items-center sm:flex-row'>
                  <div className='flex-grow mt-6 sm:pl-6 sm:mt-0'>
                    <h2 className='mb-1 text-xl font-medium text-purple-700 title-font'>
                      {t('steps:step2', 'Fill in form')}
                    </h2>
                    <p className='leading-relaxed'>
                      {t('steps:desc2', 'Enter the title and the description of the assignment. Attach files if needed.')}
                    </p>
                  </div>
                </div>
              </div>
              <div className='relative flex pb-20 mx-auto sm:items-center md:w-2/3'>
                <div className='absolute inset-0 flex items-center justify-center w-6 h-full'>
                  <div className='w-1 h-full bg-purple-200 pointer-events-none'></div>
                </div>
                <div className='relative z-10 inline-flex items-center justify-center flex-shrink-0 w-6 h-6 mt-10 text-sm font-medium text-purple-100 bg-purple-800 rounded-full sm:mt-0 title-font'>
                  3
                </div>
                <div className='flex flex-col items-start flex-grow pl-6 md:pl-8 sm:items-center sm:flex-row'>
                  <div className='flex-grow mt-6 sm:pl-6 sm:mt-0'>
                    <h2 className='mb-1 text-xl font-medium text-purple-700 title-font'>
                      {t('steps:step3', 'Generate Link')}
                    </h2>
                    <p className='leading-relaxed'>
                      {t('steps:desc3', 'Click on the button to generate a link for the homework (quiz or attendance) page and new page will be created.')}
                    </p>
                  </div>
                </div>
              </div>
              <div className='relative flex pb-10 mx-auto sm:items-center md:w-2/3'>
                <div className='absolute inset-0 flex items-center justify-center w-6 h-full'>
                  <div className='w-1 h-full bg-purple-200 pointer-events-none'></div>
                </div>
                <div className='relative z-10 inline-flex items-center justify-center flex-shrink-0 w-6 h-6 mt-10 text-sm font-medium text-purple-100 bg-purple-800 rounded-full sm:mt-0 title-font'>
                  4
                </div>
                <div className='flex flex-col items-start flex-grow pl-6 md:pl-8 sm:items-center sm:flex-row'>
                  <div className='flex-grow mt-6 sm:pl-6 sm:mt-0'>
                    <h2 className='mb-1 text-xl font-medium text-purple-700 title-font'>
                      {t('steps:step4', 'Share Link')}
                    </h2>
                    <p className='leading-relaxed'>
                      {t('steps:desc4', 'Copy the generated link and share it with your students. Track their submissions.')}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
      </div>
    )
  }
  
  export default Steps
  