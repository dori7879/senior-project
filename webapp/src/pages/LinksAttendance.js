import { useRef, useState } from 'react'
import { useTranslation } from 'react-i18next';
import Footer from '../components/Footer'
import Header from '../components/Header'


const LinksAttendance = (props) => {
  const { t } = useTranslation(['translation', 'links']);
  const [studentLinkCopied, setStudentLinkCopied] = useState(false)
  const [teacherLinkCopied, setTeacherLinkCopied] = useState(false)

  const studentTextArea = useRef(null)
  const teacherTextArea = useRef(null)

  const fullStudentLink = `/attendances/submit/${props.location.state.studentLink}`
  const fullTeacherLink = `/attendances/list-submissions/${props.location.state.teacherLink}`

  const copyLinkClipboard = (e, caller) => {
    if (caller === 'student') {
      studentTextArea.current.select()
      document.execCommand('copy')
      e.target.focus()
      setStudentLinkCopied(true)
      return
    }

    teacherTextArea.current.select()
    document.execCommand('copy')
    e.target.focus()
    setTeacherLinkCopied(true)
  }

  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center pb-4'>
          <div className='pt-4 text-xl font-bold text-purple-900'>{t('links:links', 'LinksAttendance')}</div>

          <div className='flex flex-row '>
            <div className='mr-20'>
              <h1 className='pt-4 text-xl font-bold text-gray-700'>
                {t('links:students', 'For students:')}
              </h1>
              {document.queryCommandSupported('copy') && (
                <div className='pt-4 text-xl font-bold text-purple-900'>
                  <button
                    className='relative flex justify-center w-full px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
                    onClick={(e) => copyLinkClipboard(e, 'student')}
                  >
                     {t('links:copy', 'Copy')}
                  </button>
                  {studentLinkCopied && t('links:copied', 'Copied') } 
                </div>
              )}
              <form>
                <textarea
                  ref={studentTextArea}
                  defaultValue={fullStudentLink}
                  rows='10'
                />
              </form>
            </div>
            <div>
              <h1 className='pt-4 text-xl font-bold text-gray-700'>
              {t('links:teacher', 'For teacher:')}
              </h1>
              {document.queryCommandSupported('copy') && (
                <div className='pt-4 text-xl font-bold text-purple-900'>
                  <button
                    className='relative flex justify-center w-full px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
                    onClick={(e) => copyLinkClipboard(e, 'teacher')}
                  >
                    {t('links:copy', 'Copy')}
                  </button>
                  {teacherLinkCopied && t('links:copied', 'Copied') }
                </div>
              )}
              <form>
                <textarea
                  ref={teacherTextArea}
                  defaultValue={fullTeacherLink}
                  rows='10'
                />
              </form>
            </div>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default LinksAttendance