import { useFieldArray } from "react-hook-form";
import { useTranslation } from 'react-i18next';
import { Controller } from "react-hook-form";

const Choices = ({ nestIndex, control, register, getValues, question }) => {
  // eslint-disable-next-line no-unused-vars
  const { t } = useTranslation(['translation', 'steps']);
  const { fields, remove, append } = useFieldArray({
    control,
    name: `Questions[${nestIndex}].Choices`,
    keyName: "keyID"
  });

  const renderInput = (idx, qType) => {
    if (qType === 'multiple') {
      return (
        <Controller
          control={control}
          name={`Questions[${nestIndex}].MultipleChoiceAnswer`}
          defaultValue={[]}
          render={(props) => (
            <input
              type='checkbox'
              className='mx-2 outline-none'
              onChange={(e) => {
                let current = getValues(`Questions[${nestIndex}].MultipleChoiceAnswer`)
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
      )
    } else if (qType === 'single') {
      return (
        <Controller
          control={control}
          name={`Questions[${nestIndex}].SingleChoiceAnswer`}
          defaultValue={idx}
          render={(props) => (
            <input
              type='radio'
              name={`Questions[${nestIndex}].SingleChoiceAnswer`}
              className='mx-2 outline-none'
              onChange={(e) => props.onChange(+e.target.value)}
              value={props.value}
            />
          )}
        />
      )
    }

    return null
  } 

  return (
    <div className='flex flex-row justify-between'>
      <div className="flex flex-col">
        <div className='pb-2 ml-4 text-xs text-gray-700'>
          Tick correct answers
        </div>

        
        {fields.map((choice, idx) => (
          <div key={choice.keyID} className='flex flex-row items-center mt-2'>
            <label className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700'>
              Option #{idx+1}
            </label>

            <Controller
              control={control}
              name={`Questions[${nestIndex}].Choices[${idx}].value`}
              defaultValue={choice.value}
              render={(props) => (
                <input
                  type="text"
                  onChange={(e) => {
                    props.onChange(e.target.value)
                  }}
                  value={props.value}
                  className='px-2 py-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                  placeholder='Enter your option'
                />
              )}
            />

            {renderInput(idx, question.Type)}
            
            <button
              type="button" 
              onClick={() => remove(idx)}
              className='w-5 h-5 text-purple-100 border border-purple-400 rounded focus:outline-none hover:bg-purple-400'
            >
              <svg viewBox='-40 0 427 427.001' fill='fillCurrent'>
                <path d='M232.398 154.703c-5.523 0-10 4.477-10 10v189c0 5.52 4.477 10 10 10 5.524 0 10-4.48 10-10v-189c0-5.523-4.476-10-10-10zm0 0M114.398 154.703c-5.523 0-10 4.477-10 10v189c0 5.52 4.477 10 10 10 5.524 0 10-4.48 10-10v-189c0-5.523-4.476-10-10-10zm0 0' />
                <path d='M28.398 127.121V373.5c0 14.563 5.34 28.238 14.668 38.05A49.246 49.246 0 0078.796 427H268a49.233 49.233 0 0035.73-15.45c9.329-9.812 14.668-23.487 14.668-38.05V127.121c18.543-4.922 30.559-22.836 28.079-41.863-2.485-19.024-18.692-33.254-37.88-33.258h-51.199V39.5a39.289 39.289 0 00-11.539-28.031A39.288 39.288 0 00217.797 0H129a39.288 39.288 0 00-28.063 11.469A39.289 39.289 0 0089.398 39.5V52H38.2C19.012 52.004 2.805 66.234.32 85.258c-2.48 19.027 9.535 36.941 28.078 41.863zM268 407H78.797c-17.098 0-30.399-14.688-30.399-33.5V128h250v245.5c0 18.813-13.3 33.5-30.398 33.5zM109.398 39.5a19.25 19.25 0 015.676-13.895A19.26 19.26 0 01129 20h88.797a19.26 19.26 0 0113.926 5.605 19.244 19.244 0 015.675 13.895V52h-128zM38.2 72h270.399c9.941 0 18 8.059 18 18s-8.059 18-18 18H38.199c-9.941 0-18-8.059-18-18s8.059-18 18-18zm0 0' />
                <path d='M173.398 154.703c-5.523 0-10 4.477-10 10v189c0 5.52 4.477 10 10 10 5.524 0 10-4.48 10-10v-189c0-5.523-4.476-10-10-10zm0 0' />
              </svg>
            </button>
          </div>
        ))}
      </div>

      <div>
        <button
          type="button"
          className='px-2 py-1 mt-1 mr-2 text-xs text-purple-900 border border-purple-400 rounded focus:outline-none hover:bg-purple-400'
          onClick={() => append({ value: "" })}
        >
          Add option
        </button>
      </div>
    </div>
  )
}
  
export default Choices
  