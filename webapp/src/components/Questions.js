import { useState } from 'react'
import { useFieldArray } from "react-hook-form";
import { useTranslation } from 'react-i18next';
import Single from './Single'
import Multiple from './Multiple'
import TrueFalse from './TrueFalse'
import Open from './Open'

const Questions = ({ control, register, getValues, setValue, errors }) => {
  // eslint-disable-next-line no-unused-vars
  const { t } = useTranslation(['translation', 'steps']);
  const [type, setType] = useState('single')
  const { fields, append, remove } = useFieldArray({
    control,
    name: "Questions",
    keyName: "keyID"
  });

  return (
    <>
      <ul>
      {fields.map((question, index) => {
        switch (type) {
          case 'single':
            return (
              <li key={question.keyID} className="mt-5">
                  <Single index={index} {...{ control, register, remove, question }} />
              </li>
            );
          case 'multiple':
            return (
              <li key={question.keyID} className="mt-5">
                  <Multiple index={index} {...{ control, register, remove, question }} />
              </li>
            );
          case 'truefalse':
            return (
              <li key={question.keyID} className="mt-5">
                  <TrueFalse index={index} {...{ control, register, remove, question }} />
              </li>
            );
          case 'open':
            return (
              <li key={question.keyID} className="mt-5">
                  <Open index={index} {...{ control, register, remove, question }} />
              </li>
            );
          default:
            return null
      }})}
      </ul>

      <div className='flex justify-center mt-5'>
        <div
          onClick={() => {
            switch (type) {
              case 'multiple':
                append({
                  Content: '',
                  Fixed: false,
                  Type: 'multiple',
                  Choices: [
                    {
                      Content: '',
                      Correct: false,
                    },
                  ]
                })
                break
              case 'single':
                append({
                  Content: '',
                  Fixed: false,
                  Type: 'single',
                  Choices: [
                    {
                      Content: '',
                    },
                  ],
                  SingleChoiceAnswer: null,
                })
                break
              case 'truefalse':
                append({
                  Content: '',
                  Fixed: false,
                  Type: 'truefalse',
                  TrueFalseAnswer: null,
                })
                break
              case 'open':
                append({
                  Content: '',
                  Fixed: false,
                  Type: 'open',
                  OpenAnswer: '',
                })
                break
              default:
                return
            }
          }}
          className='flex flex-col items-center justify-center w-3/4 px-4 py-5 mb-3 text-xl text-purple-700 border border-purple-800 border-dashed rounded items hover:bg-purple-800 hover:text-purple-100'
        >
          + Add Question
          <h3 className='mt-4'>
            Choose type of the question
            <select
              className='ml-4 text-gray-800 bg-purple-100 border border-purple-300'
              name='question'
              id='question'
              onChange={(e) => setType(e.target.value)}
              onClick={(e) => e.stopPropagation()}
              value={type}
            >
              <option value='single'>One answer</option>
              <option value='multiple'>Multiple answers</option>
              <option value='truefalse'>True/False</option>
              <option value='open'>Open</option>
            </select>
          </h3>
        </div>
      </div>
    </>
  )
}

export default Questions
