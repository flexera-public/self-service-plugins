require 'spec_helper'

describe 'Functional specs' , focus: true do

  def app
    Praxis::Application.instance
  end

  context 'index' do
    context 'with an incorrect response_content_type param' do
      it 'fails to validate the response' do
        get '/clouds/1/instances?response_content_type=somejunk&api_version=1.0'
        response = JSON.parse(last_response.body)
        expect(response['name']).to eq('Praxis::Exceptions::Validation')
        expect(response["message"]).to match(/Bad Content-Type/)
      end

      context 'with response validation disabled' do
        let(:praxis_config) { double('praxis_config', validate_responses: false) }
        let(:config) { double('config', praxis: praxis_config) }

        before do
          expect(Praxis::Application.instance.config).to receive(:praxis).and_return(praxis_config)
        end

        it 'fails to validate the response' do
          expect {
            get '/clouds/1/instances?response_content_type=somejunk&api_version=1.0'
          }.to_not raise_error
        end

      end
    end

  end

