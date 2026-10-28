import React, { useState } from "react";
import { View, Text, TouchableOpacity, Image, ScrollView, Modal } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";

import styled from "./styled";

export default ({ navigation }) => {
  const [modalVisivel, setModalVisivel] = useState(false);
  const urlImagem = "https://www.infoescola.com/wp-content/uploads/2018/03/sagui-de-tufos-brancos_790108426-1000x667.jpg";

  return (
    <SafeAreaView style={styled.safeArea}>

      <View style={styled.containerCabecalho}>
        <TouchableOpacity onPress={() => navigation?.goBack()}>
          <Ionicons name="arrow-back" size={28} color="#FFFFFF" />
        </TouchableOpacity>
        <Text style={styled.tituloCabecalho}>Informações do animal</Text>
      </View>

      <ScrollView style={styled.suporteCartao} showsVerticalScrollIndicator={false}>
        <View style={styled.containerPrincipal}>
          
          <View style={styled.blocoSuperior}>

            <TouchableOpacity onPress={() => setModalVisivel(true)} activeOpacity={0.8}>
              <Image
                source={{ uri: urlImagem }}
                style={styled.imagemAnimal}
              />
            </TouchableOpacity>

            <View style={styled.divisorVertical} />

            <View style={styled.containerInformacoes}>
              <Text style={styled.nomeAnimal}>SAGUI-DE-TUFOS-BRANCOS</Text>
              <View style={styled.divisorHorizontal} />
              <Text style={styled.nomeCientifico}>(Callithrix jacchus)</Text>
            
              <View style={styled.linhaInformacao}>
                <Text style={styled.textoInformacao}>Peso médio:</Text>
                <Text style={styled.valorInformacao}>300 g - 450 g</Text>
              </View>

            {/* Tamanho animal */}
              <View style={styled.linhaInformacao}>
                <Text style={styled.textoInformacao}>Tamanho médio:</Text>
                <Text style={styled.valorInformacao}>18 cm - 25 cm</Text>
              </View>
            </View>
          </View>

          {/* Bloco de Descrição */}
          <View style={styled.blocoDescricao}>
            <Text style={styled.tituloDescricao2}>
              Pequeno primata arbóreo com distintos tufos brancos nas orelhas.
            </Text>

            <View style={styled.linhaInformacao2}>
              <Text style={styled.textoInformacao2}>Hábito</Text>
              <Text style={styled.valorInformacao2}>Diurno</Text>
            </View>

            <View style={styled.linhaInformacao2}>
              <Text style={styled.textoInformacao2}>Conservação</Text>
              <Text style={styled.valorInformacao2}>Pouco preocupante</Text>
            </View>

            <View style={styled.linhaInformacao2}>
              <Text style={styled.textoInformacao2}>Raridade</Text>
              <Text style={styled.valorInformacao2}>2/5</Text>
            </View>

            <Text style={styled.tituloDescricao}>Vocalização / Postura: </Text> 
            <Text style={styled.valorInformacao2}>
              Emite piados agudos e ergue os tufos em sinal de alerta.
            </Text>

            {/* Container de Aviso */}
            <View style={styled.containerAviso}>
              <Ionicons name="warning" size={20} color="#D84315" style={styled.iconeAviso} />
              <Text style={styled.tituloAviso}>Aviso de Segurança: </Text>
              <Text style={styled.textoAviso}>
                Nunca tente alimentar, tocar ou capturar os animais, pois há risco de mordidas e transmissão de zoonoses.
              </Text>

              <Text style={styled.tituloAviso}>Protocolo: </Text>
              <Text style={styled.textoAviso}>
                Manter distância segura, não oferecer comida e apenas observar sem interagir.
              </Text>

              <View style={styled.cabecalhoAviso}>
                <Text style={styled.tituloAviso}>Resgate: </Text>
                <Text style={styled.textoAviso}>193 / Setor Ambiental</Text>
              </View>
            </View>

          </View>
        </View>


      </ScrollView>
        {/* Logo Discovery Rural */}
        <Image
          style={{ width: 96, height: 96, alignSelf: "center", marginTop: 15, marginBottom: 15 }}
          source={{ uri: "https://github.com/rafael9577/discoveryRural/blob/main/front-end/assets/logoDiscoveryRuralWhite.png?raw=true" }}
        />

      {/* Modal de Imagem Expandida */}
      <Modal
        visible={modalVisivel}
        transparent={true}
        animationType="fade"
        onRequestClose={() => setModalVisivel(false)} // Trata o botão Voltar do celular
      >
        <TouchableOpacity
          style={styled.modalContainer}
          activeOpacity={1}
          onPress={() => setModalVisivel(false)} // Fecha ao tocar fora/na tela
        >
          <TouchableOpacity
            style={styled.botaoFecharModal}
            onPress={() => setModalVisivel(false)}
          >
            <Ionicons name="close" size={32} color="#FFFFFF" />
          </TouchableOpacity>

          <Image
            source={{ uri: urlImagem }}
            style={styled.imagemModal}
            resizeMode="contain"
          />
        </TouchableOpacity>
      </Modal>
    </SafeAreaView>
  );
};