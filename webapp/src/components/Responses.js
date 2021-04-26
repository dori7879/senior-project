import { useTranslation } from 'react-i18next';
import { Controller } from "react-hook-form";

const Responses = ({ control, register, getValues, setValue, errors, questions }) => {
  // eslint-disable-next-line no-unused-vars
  const { t } = useTranslation(['translation', 'steps']);

  return (
    <>
      <ul>
      {questions.map((question, index) => {
        switch (question.Type) {
          case 1:
            return (
              <li className="mt-5" key={question.ID}>
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
                          <Controller
                            control={control}
                            name={`Responses[${index}].SingleChoiceResponse`}
                            defaultValue={idx}
                            render={(props) => (
                              <input
                                type='radio'
                                name={`Responses[${index}].SingleChoiceResponse`}
                                className='mx-2 outline-none'
                                onChange={(e) => props.onChange(+e.target.value)}
                                value={props.value}
                              />
                            )}
                          />
                          <div className='px-2 py-1 text-xs leading-tight text-gray-700 focus:outline-none focus:bg-white'>
                            {choice}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                  <Controller
                    control={control}
                    name={`Responses[${index}].QuestionID`}
                    defaultValue={question.ID}
                    render={(props) => (
                      <input
                        type="hidden"
                        name={`Responses[${index}].QuestionID`}
                        className='outline-none'
                        onChange={(e) => props.onChange(+e.target.value)}
                        value={props.value}
                      />
                    )}
                  />
                  <Controller
                    control={control}
                    name={`Responses[${index}].Type`}
                    defaultValue={1}
                    render={(props) => (
                      <input
                        type="hidden"
                        name={`Responses[${index}].Type`}
                        className='outline-none'
                        onChange={(e) => props.onChange(+e.target.value)}
                        value={props.value}
                      />
                    )}
                  />
                </div>
              </li>
            );
          case 2:
            return (
              <li className="mt-5" key={question.ID}>
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
                          <Controller
                            control={control}
                            name={`Responses[${index}].MultipleChoiceResponse`}
                            defaultValue={[]}
                            render={(props) => (
                              <input
                                type='checkbox'
                                className='mx-2 outline-none'
                                onChange={(e) => {
                                  let current = getValues(`Responses[${index}].MultipleChoiceResponse`)
                                  if (current.includes(+e.target.value)) {
                                    props.onChange(current.filter(el => el !== +e.target.value))
                                  } else {
                                    props.onChange([...current, +e.target.value])
                                  }
                                }}
                                value={idx}
                              />
                            )}
                          />
                          <div className='px-2 py-1 text-xs leading-tight text-gray-700 focus:outline-none focus:bg-white'>
                            {choice}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                  <Controller
                    control={control}
                    name={`Responses[${index}].QuestionID`}
                    defaultValue={question.ID}
                    render={(props) => (
                      <input
                        type="hidden"
                        name={`Responses[${index}].QuestionID`}
                        className='outline-none'
                        onChange={(e) => props.onChange(+e.target.value)}
                        value={props.value}
                      />
                    )}
                  />
                  <Controller
                    control={control}
                    name={`Responses[${index}].Type`}
                    defaultValue={2}
                    render={(props) => (
                      <input
                        type="hidden"
                        name={`Responses[${index}].Type`}
                        className='outline-none'
                        onChange={(e) => props.onChange(+e.target.value)}
                        value={props.value}
                      />
                    )}
                  />
                </div>
              </li>
            );
          case 3:
            return (
              <li className="mt-5" key={question.ID}>
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
                    <Controller
                      control={control}
                      name={`Responses[${index}].TrueFalseResponse`}
                      defaultValue={true}
                      render={(props) => (
                        <input
                          type="radio"
                          name={`Responses[${index}].TrueFalseResponse`}
                          className='outline-none'
                          onChange={(e) => props.onChange(e.target.value === 'true')}
                          value={props.value.toString()}
                        />
                      )}
                    />
                    <label
                      htmlFor='choice'
                      className='px-1 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase '
                    >
                      true
                    </label>
                    <br></br>
                    <Controller
                      control={control}
                      name={`Responses[${index}].TrueFalseResponse`}
                      defaultValue={false}
                      render={(props) => (
                        <input
                          type="radio"
                          name={`Responses[${index}].TrueFalseResponse`}
                          className='outline-none'
                          onChange={(e) => props.onChange(e.target.value === 'true')}
                          value={props.value.toString()}
                        />
                      )}
                    />
                    <label
                      htmlFor='choice'
                      className='px-1 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase '
                    >
                      false
                    </label>
                    <br></br>
                  </div>
                  <Controller
                    control={control}
                    name={`Responses[${index}].QuestionID`}
                    defaultValue={question.ID}
                    render={(props) => (
                      <input
                        type="hidden"
                        name={`Responses[${index}].QuestionID`}
                        className='outline-none'
                        onChange={(e) => props.onChange(+e.target.value)}
                        value={props.value}
                      />
                    )}
                  />
                  <Controller
                    control={control}
                    name={`Responses[${index}].Type`}
                    defaultValue={3}
                    render={(props) => (
                      <input
                        type="hidden"
                        name={`Responses[${index}].Type`}
                        className='outline-none'
                        onChange={(e) => props.onChange(+e.target.value)}
                        value={props.value}
                      />
                    )}
                  />
                </div>
              </li>
            );
          case 4:
            return (
              <li className="mt-5" key={question.ID}>
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
                  <Controller
                    control={control}
                    name={`Responses[${index}].QuestionID`}
                    defaultValue={question.ID}
                    render={(props) => (
                      <input
                        type="hidden"
                        name={`Responses[${index}].QuestionID`}
                        className='outline-none'
                        onChange={(e) => props.onChange(+e.target.value)}
                        value={props.value}
                      />
                    )}
                  />
                  <Controller
                    control={control}
                    name={`Responses[${index}].Type`}
                    defaultValue={4}
                    render={(props) => (
                      <input
                        type="hidden"
                        name={`Responses[${index}].Type`}
                        className='outline-none'
                        onChange={(e) => props.onChange(+e.target.value)}
                        value={props.value}
                      />
                    )}
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
  