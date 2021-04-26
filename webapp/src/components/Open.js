import { useTranslation } from 'react-i18next';

const Open = ({ index, control, register, remove, question }) => {
    // eslint-disable-next-line no-unused-vars
    const { t } = useTranslation(['translation', 'steps']);
    
    return (
      <div className='flex flex-col'>
        <div className='flex flex-row'>
          <div className='flex flex-col w-full'>
            <div className='flex flex-row justify-between w-full px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
              <div className='flex flex-row'>
                <label className='block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-purple-800'>
                  Question #{index+1}
                </label>

                <input
                  name={`Questions[${index}].Content`}
                  defaultValue={question.Content} // make sure to set up defaultValue
                  ref={register()}
                  className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded w-80 focus:outline-none focus:bg-white'
                  placeholder='Enter your question'
                />

                <input
                  ref={register()}
                  control={control}
                  type='hidden'
                  name={`Questions[${index}].Type`}
                  className='outline-none'
                  defaultValue={'open'}
                />
              </div>

              <div>
                <button
                  onClick={() => remove(index)}
                  className='w-5 h-5 text-purple-100 border border-purple-400 rounded focus:outline-none hover:bg-purple-400'
                >
                  <svg
                    xmlns='http://www.w3.org/2000/svg'
                    viewBox='0 0 92 92'
                    fill='#5B21B6'
                  >
                    <path d='M70.7 64.3c1.8 1.8 1.8 4.6 0 6.4-.9.9-2 1.3-3.2 1.3-1.2 0-2.3-.4-3.2-1.3L46 52.4 27.7 70.7c-.9.9-2 1.3-3.2 1.3s-2.3-.4-3.2-1.3a4.47 4.47 0 010-6.4L39.6 46 21.3 27.7a4.47 4.47 0 010-6.4c1.8-1.8 4.6-1.8 6.4 0L46 39.6l18.3-18.3c1.8-1.8 4.6-1.8 6.4 0 1.8 1.8 1.8 4.6 0 6.4L52.4 46l18.3 18.3z' />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <div className='flex flex-row'>
            <label className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                Answer
            </label>
            <textarea
                defaultValue={question.OpenAnswer}
                name={`Questions[${index}].OpenAnswer`}
                ref={register()}
                className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                placeholder='Enter your answer'
            />
        </div>
      </div>
    )
  }
  
  export default Open
  