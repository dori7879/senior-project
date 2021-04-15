import { useTranslation } from 'react-i18next';

const Responses = ({ control, register, getValues, setValue, errors, questions }) => {
  // eslint-disable-next-line no-unused-vars
  const { t } = useTranslation(['translation', 'steps']);

  return (
    <>
      <ul>
      {questions.map((question, index) => {
        switch (question.Type) {
          case 'single':
            return (
              <li key={question.ID}>
                <div className='flex flex-col'>
                  <div className='flex flex-row'>
                    <div className='flex flex-col w-full'>
                      <div className='flex flex-row justify-between w-full px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                        <div className='flex flex-row'>
                          <label className='block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-purple-800'>
                            Question
                          </label>
                          <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded w-80 focus:outline-none focus:bg-white'>
                            {question.Content}
                          </div>
                        </div>
                      </div>
                      <div className='pb-2 ml-4 text-xs text-gray-700'>
                        Tick the correct answer
                      </div>
                      {question.Choices.map((choice, idx) => (
                        // CHANGE KEY FROM idx
                        <div key={idx} className='flex flex-row'>
                          <input
                            ref={register()}
                            control={control}
                            type='radio'
                            name={`Responses[${index}].SingleChoiceResponse`}
                            defaultValue={idx} // make sure to set up defaultValue
                            className='mx-2 outline-none'
                          />
                          <div className='px-2 py-1 text-xs leading-tight text-gray-700 focus:outline-none focus:bg-white'>
                            {choice}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                  <input
                    ref={register()}
                    control={control}
                    type='hidden'
                    name={`Responses[${index}].QuestionID`}
                    className='outline-none'
                    defaultValue={question.ID}
                  />
                </div>
              </li>
            );
          case 'multiple':
            return (
              <li key={question.ID}>
                <div className='flex flex-col'>
                  <div className='flex flex-row'>
                    <div className='flex flex-col w-full'>
                      <div className='flex flex-row justify-between w-full px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                        <div className='flex flex-row'>
                          <label className='block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-purple-800'>
                            Question
                          </label>
                          <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded w-80 focus:outline-none focus:bg-white'>
                            {question.Content}
                          </div>
                        </div>
                      </div>
                      <div className='pb-2 ml-4 text-xs text-gray-700'>
                        Tick correct answers
                      </div>
                      {question.Choices.map((choice, idx) => (
                        <div key={idx} className='flex flex-row'>
                          <input
                            ref={register()}
                            control={control}
                            type='checkbox'
                            name={`Responses[${index}].MultipleChoiceResponse`}
                            defaultValue={idx} // make sure to set up defaultValue
                            className='mx-2 outline-none'
                          />
                          <div className='px-2 py-1 text-xs leading-tight text-gray-700 focus:outline-none focus:bg-white'>
                            {choice}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                  <input
                    ref={register()}
                    control={control}
                    type='hidden'
                    name={`Responses[${index}].QuestionID`}
                    className='outline-none'
                    defaultValue={question.ID}
                  />
                </div>
              </li>
            );
          case 'truefalse':
            return (
              <li key={question.ID}>
                <div className='flex flex-col w-full'>
                  <div className='flex flex-row justify-between px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                    <div className='flex flex-row'>
                      <label className='block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-purple-800'>
                        Question
                      </label>
                      <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded w-80 focus:outline-none focus:bg-white'>
                        {question.Content}
                      </div>
                    </div>
                  </div>
                  <div className='ml-4'>
                    <input
                      ref={register()}
                      control={control}
                      type='radio'
                      name={`Responses[${index}].TrueFalseResponse`}
                      className='outline-none'
                      defaultValue={true}
                    />
                    <label
                      htmlFor='choice'
                      className='px-1 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase '
                    >
                      true
                    </label>
                    <br></br>
                    <input
                      ref={register()}
                      control={control}
                      type='radio'
                      name={`Responses[${index}].TrueFalseResponse`}
                      className='outline-none'
                      defaultValue={false}
                    />
                    <label
                      htmlFor='choice'
                      className='px-1 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase '
                    >
                      false
                    </label>
                    <br></br>
                  </div>
                  <input
                    ref={register()}
                    control={control}
                    type='hidden'
                    name={`Responses[${index}].QuestionID`}
                    className='outline-none'
                    defaultValue={question.ID}
                  />
                </div>
              </li>
            );
          case 'open':
            return (
              <li key={question.ID}>
                <div className='flex flex-col w-full'>
                  <div className='flex flex-row justify-between px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                    <div className='flex flex-row'>
                      <label className='block px-2 pt-1 mb-2 text-xs font-bold tracking-wide text-purple-800'>
                        Question
                      </label>
                      <div className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded w-80 focus:outline-none focus:bg-white'>
                        {question.Content}
                      </div>
                    </div>
                  </div>
                  <div className='flex flex-row'>
                    <label className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
                      Answer
                    </label>
                    <textarea
                      ref={register()}
                      control={control}
                      name={`Responses[${index}].OpenResponse`}
                      defaultValue={""}
                      className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                      placeholder='Enter your answer'
                    />
                  </div>
                  <input
                    ref={register()}
                    control={control}
                    type='hidden'
                    name={`Responses[${index}].QuestionID`}
                    className='outline-none'
                    defaultValue={question.ID}
                  />
                </div>
              </li>
            );
          default:
            return null
      }})}
      </ul>
    </>
  )
}

export default Responses
  