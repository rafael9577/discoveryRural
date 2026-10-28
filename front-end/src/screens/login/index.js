import { View, Text, TouchableOpacity, TextInput, Image, Modal } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import Toast from "react-native-toast-message";
import { Feather } from "@expo/vector-icons"
import { useState } from 'react';
import styled from "./styled";

import * as Api from "../../../services/API"

export default () => {

    const [imagePassowrd, setImagePassowrd] = useState("eye-off");
    const [passowrdFild, setPassowrdFild] = useState("");
    const [emailFild, setEmailFild] = useState("");
    const [seePassowrd, setSeePassowrd] = useState(true);
    const [loading, setLoading] = useState(false);

    const Navigator = useNavigation();

    const tosignUp = () => {
        Navigator.navigate("signup");
    }

    async function handleLoginButtonClick() {
        if (emailFild != '' && passowrdFild != '') {
            setLoading(true)
            try {
                let reqLogin = await Api.UserSingIn(emailFild.trim().toLocaleLowerCase(), passowrdFild.trim())
                if (Object.keys(reqLogin)[0] == 'error') {
                    Toast.show({
                        type: 'error',
                        text1: 'Email ou senha invalido'
                    })
                } else {
                    console.log(reqLogin)
                    Navigator.replace("tabsRouts");
                }
            } catch (error) {
                console.log(error)
                Toast.show({
                    type: 'error',
                    text1: error.message
                })
            } finally {
                setLoading(false)
            }
        } else {
            Toast.show({
                type: 'info',
                text1: "Preencha todos os campos",
            });
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
                <Image style= {{ width: 128, height:128, alignSelf: 'center'}}source={{uri: "https://github.com/rafael9577/discoveryRural/blob/main/front-end/assets/logoDiscoveryRuralWhite.png?raw=true"}} /> 
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
                    <TouchableOpacity  style={{ marginRight: 12 }} onPress={visibiliyPassowrd}>
                        <Feather name={imagePassowrd} size={20} color="black"/>
                    </TouchableOpacity>
                </View>
                
                <Text style={styled.forgotPassowrd}>Esqueci a senha</Text>
                
                <TouchableOpacity style={styled.custtomButton} onPress={handleLoginButtonClick}>
                    <Text style={styled.customButtontext}>Login</Text>
                </TouchableOpacity>
            </View>
        
            <View  >
                <TouchableOpacity onPress={tosignUp} style={styled.singMessageButto}>
                    <Text > Ainda não possui uma conta? </Text>
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