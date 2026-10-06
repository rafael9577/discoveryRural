import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, Modal } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import Toast from "react-native-toast-message";
import { Feather } from "@expo/vector-icons"
import { useState } from 'react';
import styled from "./styled";


export default () => {

    const [imagePassowrd, setImagePassowrd] = useState("eye-off");
    const [passowrdFild, setPassowrdFild] = useState("123");
    const [emailFild, setEmailFild] = useState('');
    const [seePassowrd, setSeePassowrd] = useState(true);
    const [loading, setLoading] = useState(false);

    const Navigator = useNavigation();

    const userNavigate = async () => {
        
    }

    async function handleLoginButtonClick() {
        if (emailFild != '' && passowrdFild != '') {
            setLoading(true)
            try{
                // chamar a API para validar o usuario
                true
            } catch (error){
                // retornar o erro que deu ao tentar validar
                Toast
            }
        }
    }

    const visibiliyPassowrd = () => {
        if (!seePassowrd) {
            setImagePassowrd("eye-off")
        } else {
            setImagePassowrd("eye")
        }
        setSeePassowrd(!seePassowrd)
    }

    return (
        <SafeAreaView style={styled.safeArea}>
            <View style= {styled.inputArea}>
                <Image style= {{ width: 128, height:128, alignSelf: 'center'}}source={require('../../../assets/logoDiscoveryRuralWhite.png')} /> 
                <View style={styled.singInPut}>
                    <TextInput 
                        style={styled.textInput}
                        placeholder="Digitel seu login/Email"
                        placeholderTextColor='#5E6737'
                        value={emailFild}
                        onChangeText={t => setEmailFild(t)}
                        />
                </View>
                
                <View style={styled.singInPut}>
                    <TextInput 
                        style={styled.textInput}
                        placeholder="Digite sua Senha"
                        placeholderTextColor='#5E6737'
                        value={passowrdFild}
                        secureTextEntry={seePassowrd}
                        onChangeText={t => setPassowrdFild(t)}
                    />
                    <TouchableOpacity  style={{ marginRight: 12 }} onPress={visibiliyPassowrd} >
                        <Feather name={imagePassowrd} size={20} color="black"/>
                    </TouchableOpacity>
                </View>
                
                <Text style={styled.forgotPassowrd}>Esqueci a senha</Text>
            </View>
        
            <View >
                <TouchableOpacity style={styled.custtomButton} >
                    <Text style={styled.customButtontext}>Login</Text>
                </TouchableOpacity>

                <TouchableOpacity onPress={''} style={styled.singMessageButto} >
                    <Text>Ainda não possui uma conta?</Text>
                    <Text style={styled.singMessageButtonTextBold}>Cadastre-se</Text>
                </TouchableOpacity>    
            </View>

                <Modal
                    transparent={true}
                    animationType="none"
                    visible={loading}
                    />

        </SafeAreaView>
    );
};