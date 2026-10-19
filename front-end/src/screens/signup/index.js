import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, Modal } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import Toast from "react-native-toast-message";
import { Feather } from "@expo/vector-icons"
import { useState } from 'react';
import styled from "./styled";

import * as Api from "../../../services/API"


export default () => {
    
    const navigation = useNavigation();

    const [loading, setloading] = useState(false);
    const [nomeFild, setNomeFild] = useState("");
    const [cpfFild, setCpfFild] = useState("");
    const [emailFild, setEmailFild] = useState("");
    const [passowrdFild, setPassowrdFild] = useState("");
    const [passowrdFild1, setPassowrdFild1] = useState("");
    const [seePass, setSeePass] = useState(true);
    const [seePass1, setSeePass1] = useState(true);
    
    const validarEmail = /^[a-z0-9.]+@[a-z0-9]+\.[a-z]/i

        function validarCPF(cpf) {
            if (cpf.length !== 11) return false;

            if (/^(\d)\1{10}$/.test(cpf)) return false;

            let soma = 0;
            for (let i = 0; i < 9; i++) {
                soma += parseInt(cpf.charAt(i)) * (10 - i);
            }
            let resto = (soma * 10) % 11;
            if (resto === 10 || resto === 11) resto = 0;
            if (resto !== parseInt(cpf.charAt(9))) return false;

            
            soma = 0;
            for (let i = 0; i < 10; i++) {
                soma += parseInt(cpf.charAt(i)) * (11 - i);
            }
            resto = (soma * 10) % 11;
            if (resto === 10 || resto === 11) resto = 0;
            if (resto !== parseInt(cpf.charAt(10))) return false;

            return true;
        }

        function validacao() {
        if (emailFild != '' && nomeFild != '' && passowrdFild != '' && cpfFild != '') {
            if ((validarEmail).test(emailFild.toLocaleLowerCase())) {
                if (passowrdFild == passowrdFild1) {
                    if (validarCPF(cpfFild)) {
                        return true
                    } else {
                        Toast.show({
                        type: 'error',
                        text1: 'CPF invalido!'
                    })
                    }
                } else {
                    Toast.show({
                        type: 'error',
                        text1: 'As senhas devem ser iguais!'
                    })
                }
            } else {
                Toast.show({
                    type: 'error',
                    text1: 'Email invalido'
                })
            }
        } else {
            Toast.show({
                type: 'error',
                text1: 'Preencha todos os campos ou digite um email valido',
            })
        }
        return false
    }

    const handleSignUp = async () => {
        if (validacao()) {
            setloading(true)
            try {
                let resp = await Api.cadastrar({
                    nome_usuario: nomeFild.trim(),
                    cpf: cpfFild.replace(/\D/g, ''), 
                    email: emailFild.toLowerCase().trim(),
                    senha: passowrdFild.trim()}
                );
                if (Object.keys(resp)[0] == 'error') {
                    Toast.show({
                        type: 'error',
                        text1: error.message
                    })
                } else {
                        navigation.reset({
                            routes: [{ name: 'login' }]
                        })
                }
            } catch (error) {
                console.log(error)
                Toast.show({
                    type: 'error',
                    text1: error.message
                })
            } finally {
                setloading(false)
            }
        }
    }

    const toLogin = () => {
        navigation.navigate("login");
    };

    return (
        <SafeAreaView style={styled.safeArea}>
            <View style={styled.inputArea}>
                <Image style= {{ width: 128, height:128, alignSelf: 'center'}}source={{uri: "https://github.com/rafael9577/discoveryRural/blob/main/front-end/assets/logoDiscoveryRuralWhite.png?raw=true"}} /> 
                
                <View style={[styled.singInPut, { flexDirection: 'row', alignItems: 'center' }]}>
                    <TextInput
                        style={styled.textInput}
                        placeholder="insira o seu nome"
                        placeholderTextColor="#666666"
                        value={nomeFild}   
                        onChangeText={t => setNomeFild(t)}
                    />
                </View>
                
                <View style={[styled.singInPut, { flexDirection: 'row', alignItems: 'center' }]}>
                    <TextInput
                        style={styled.textInput}
                        placeholder="insira o CPF"
                        placeholderTextColor="#666666"
                        value={cpfFild}   
                        onChangeText={t => setCpfFild(t)}
                    />
                </View>

                <View style={styled.singInPut}>
                    <TextInput
                        style={styled.textInput}
                        placeholder="insira seu e-mail"
                        placeholderTextColor="#666666"
                        value={emailFild}
                        onChangeText={t => setEmailFild(t)}
                    />
                </View>


                <View style={[styled.singInPut, { flexDirection: 'row', alignItems: 'center' }]}>
                    <TextInput
                        style={styled.textInput}
                        placeholder="insira a senha"
                        placeholderTextColor="#666666"
                        secureTextEntry={seePass} // para ocultar os caracteres
                        value={passowrdFild}   //
                        onChangeText={t => setPassowrdFild(t)}
                    />
                    <TouchableOpacity onPress={() => setSeePass(!seePass)} style={{ marginRight: 12 }}>
                        <Feather 
                            name={seePass ? "eye-off" : "eye"} 
                            size={20} 
                            color="black" 
                        />
                    </TouchableOpacity>
                </View>

                <View style={[styled.singInPut, { flexDirection: 'row', alignItems: 'center' }]}>
                    <TextInput
                        style={styled.textInput}
                        placeholder="confirme a senha"
                        placeholderTextColor="#666666"
                        secureTextEntry={seePass1} 
                        value={passowrdFild1}   
                        onChangeText={t => setPassowrdFild1(t)}
                    />
                    <TouchableOpacity onPress={() => setSeePass1(!seePass1)} style={{ marginRight: 12 }}>
                        <Feather 
                            name={seePass1 ? "eye-off" : "eye"} 
                            size={20} 
                            color="black" 
                        />
                    </TouchableOpacity>
                </View>
            </View>
            
            <TouchableOpacity style={styled.custtomButton} onPress={handleSignUp}>
                <Text style={styled.customButtontext}>Cadastrar</Text>
            </TouchableOpacity>  

            <TouchableOpacity onPress={toLogin} style={styled.singMessageButto}>
                <Text style={styled.singMessageButtonTextBold}> Já possui uma conta? </Text>
            </TouchableOpacity>
            
            {loading && (
                <View pointerEvents="none" style={styled.loadingPoint} >
                    <ActivityIndicator size='large' color="#0000ff" />
                </View>
            )}

            <Modal
                transparent={true}
                animationType="none"
                visible={loading}
            />
        </SafeAreaView>
    );
}
